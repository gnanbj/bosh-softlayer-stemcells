package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

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
		Context("SUCCESSFULLY clean up the VGBDTG", func() {
			Context("when using name pattern matching", func() {
				It("issues clean-up command and verifies VGBDTG is deleted", func() {

				})
			})

			Context("when using name pattern matching and date", func() {

			})
		})

		Context("FAILS to cleans up the new VGBDTG", func() {
			Context("when name pattern doesn't match", func() {

			})

			Context("when date provided is in the past", func() {

			})

			Context("when name pattern matches and date is in the future, but is tagged with SHIPIT", func() {

			})
		})
	})
})
