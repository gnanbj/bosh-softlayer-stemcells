package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	integrationHelpers "github.com/maximilien/bosh-softlayer-stemcells/test_helpers"
	data_types "github.com/maximilien/softlayer-go/data_types"
	softlayer "github.com/maximilien/softlayer-go/softlayer"
	testhelpers "github.com/maximilien/softlayer-go/test_helpers"
)

var _ = FDescribe("Cleanup Stemcells - Integration", func() {
	var (
		err           error
		vgbdtgService softlayer.SoftLayer_Virtual_Guest_Block_Device_Template_Group_Service
	)

	BeforeEach(func() {
		vgbdtgService, err = testhelpers.CreateVirtualGuestBlockDeviceTemplateGroupService()
		Expect(err).ToNot(HaveOccurred())
	})

	Context("creates a VGBDTG with a test name", func() {
		BeforeEach(func() {
			//Creating VGBDTG from service using test VHD file
			swift_username := strings.Split(os.Getenv("SWIFT_USERNAME"), ":")[0]
			if swift_username == "" {
				fmt.Println("SWIFT_USERNAME environment must be set")
			}

			swift_cluster := os.Getenv("SWIFT_CLUSTER")
			if swift_username == "" {
				fmt.Println("SWIFT_USERNAME environment must be set")
			}

			vgbdtgConfig := data_types.SoftLayer_Container_Virtual_Guest_Block_Device_Template_Configuration{
				Name: "integration-test-vgbtg",
				Note: "",
				OperatingSystemReferenceCode: "UBUNTU_14_64",
				Uri: "swift://" + swift_username + "@" + swift_cluster + "/stemcells/test-bosh-stemcell-softlayer.vhd",
			}

			_, err := vgbdtgService.CreateFromExternalSource(vgbdtgConfig)
			Expect(err).ToNot(HaveOccurred())

			foundVGBTG := searchForVGBDTG("integration-test-vgbtg")
			Expect(foundVGBTG).To(Equal(true))

			//Wait for transaction to complete
			time.Sleep(1 * time.Minute)
		})

		AfterEach(func() {
			//Make sure NO VGBTG objects exist with name integration-test-vgbtg
			pwd, err := os.Getwd()
			Expect(err).ToNot(HaveOccurred())

			rootPath := filepath.Join(pwd, "..")

			_, err = integrationHelpers.RunSLStemcells(rootPath, []string{"-c", "cleanup-stemcells", "--name-pattern", "integration-test-vgbtg", "--last-valid-date", time.Now().AddDate(1, 0, 0).Format("2006-01-02"), "--shipit-tag", "fake-tag"})
			Expect(err).ToNot(HaveOccurred())

			//Wait for transaction to complete
			time.Sleep(1 * time.Minute)

			foundVGBTG := searchForVGBDTG("integration-test-vgbtg")
			Expect(foundVGBTG).To(Equal(false))
		})

		Context("SUCCESSFULLY clean up the VGBDTG", func() {
			Context("when using name pattern matching and date is in the future", func() {
				It("issues clean-up command and verifies VGBDTG is deleted", func() {
					pwd, err := os.Getwd()
					Expect(err).ToNot(HaveOccurred())

					rootPath := filepath.Join(pwd, "..")

					_, err = integrationHelpers.RunSLStemcells(rootPath, []string{"-c", "cleanup-stemcells", "--name-pattern", "integration-test-vgbtg", "--last-valid-date", time.Now().AddDate(1, 0, 0).Format("2006-01-02")})
					Expect(err).ToNot(HaveOccurred())

					//Wait for transaction to complete
					time.Sleep(1 * time.Minute)

					foundVGBTG := searchForVGBDTG("integration-test-vgbtg")
					Expect(foundVGBTG).To(Equal(false))
				})
			})
		})

		Context("FAILS to cleans up the new VGBDTG", func() {
			Context("when name pattern doesn't match", func() {
				It("issues a clean-up command and verifies VGBDTG is NOT deleted", func() {
					pwd, err := os.Getwd()
					Expect(err).ToNot(HaveOccurred())

					rootPath := filepath.Join(pwd, "..")

					_, err = integrationHelpers.RunSLStemcells(rootPath, []string{"-c", "cleanup-stemcells", "--name-pattern", "fake-pattern", "--last-valid-date", time.Now().AddDate(1, 0, 0).Format("2006-01-02")})
					Expect(err).ToNot(HaveOccurred())

					//Wait for transaction to complete
					time.Sleep(1 * time.Minute)

					foundVGBTG := searchForVGBDTG("integration-test-vgbtg")
					Expect(foundVGBTG).To(Equal(true))
				})
			})

			Context("when date provided is in the past", func() {
				It("issues a clean-up command and verifies VGBDTG is NOT deleted", func() {
					pwd, err := os.Getwd()
					Expect(err).ToNot(HaveOccurred())

					rootPath := filepath.Join(pwd, "..")

					_, err = integrationHelpers.RunSLStemcells(rootPath, []string{"-c", "cleanup-stemcells", "--name-pattern", "integration-test-vgbtg", "--last-valid-date", "2010-09-23"})
					Expect(err).ToNot(HaveOccurred())

					//Wait for transaction to complete
					time.Sleep(1 * time.Minute)

					foundVGBTG := searchForVGBDTG("integration-test-vgbtg")
					Expect(foundVGBTG).To(Equal(true))
				})

			})

			Context("when name pattern matches and date is in the future, but is tagged with SHIPIT", func() {
				BeforeEach(func() {
					//Creating VGBDTG from service using test VHD file
					swift_username := strings.Split(os.Getenv("SWIFT_USERNAME"), ":")[0]
					if swift_username == "" {
						fmt.Println("SWIFT_USERNAME environment must be set")
					}

					swift_cluster := os.Getenv("SWIFT_CLUSTER")
					if swift_username == "" {
						fmt.Println("SWIFT_USERNAME environment must be set")
					}

					vgbdtgConfig := data_types.SoftLayer_Container_Virtual_Guest_Block_Device_Template_Configuration{
						Name: "integration-test-vgbtg",
						Note: "SHIPIT",
						OperatingSystemReferenceCode: "UBUNTU_14_64",
						Uri: "swift://" + swift_username + "@" + swift_cluster + "/stemcells/test-bosh-stemcell-softlayer.vhd",
					}

					_, err := vgbdtgService.CreateFromExternalSource(vgbdtgConfig)
					Expect(err).ToNot(HaveOccurred())

					foundVGBTG := searchForVGBDTG("integration-test-vgbtg")
					Expect(foundVGBTG).To(Equal(true))

					//Wait for transaction to complete
					time.Sleep(1 * time.Minute)
				})

				It("issues a clean-up command and verifies VGBDTG is NOT deleted", func() {
					pwd, err := os.Getwd()
					Expect(err).ToNot(HaveOccurred())

					rootPath := filepath.Join(pwd, "..")

					_, err = integrationHelpers.RunSLStemcells(rootPath, []string{"-c", "cleanup-stemcells", "--name-pattern", "integration-test-vgbtg"})
					Expect(err).ToNot(HaveOccurred())

					//Wait for transaction to complete
					time.Sleep(1 * time.Minute)

					foundVGBTG := searchForVGBDTG("integration-test-vgbtg")
					Expect(foundVGBTG).To(Equal(true))
				})
			})
		})
	})
})

func searchForVGBDTG(name string) bool {
	accountService, err := testhelpers.CreateAccountService()
	Expect(err).ToNot(HaveOccurred())

	vgbdtgObjects, err := accountService.GetBlockDeviceTemplateGroups()
	Expect(err).ToNot(HaveOccurred())

	compiledPattern, err := regexp.Compile(name)
	Expect(err).ToNot(HaveOccurred())

	filteredObjects := []data_types.SoftLayer_Virtual_Guest_Block_Device_Template_Group{}
	for _, object := range vgbdtgObjects {
		if compiledPattern.MatchString(object.Name) {
			filteredObjects = append(filteredObjects, object)
		}
	}

	if len(filteredObjects) > 0 {
		return true
	}
	return false
}
