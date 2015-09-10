package cleanup_stemcells_test

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/maximilien/bosh-softlayer-stemcells/common"

	slclientfakes "github.com/maximilien/softlayer-go/client/fakes"

	cmds "github.com/maximilien/bosh-softlayer-stemcells/cmds"
	cleanup_stemcells "github.com/maximilien/bosh-softlayer-stemcells/cmds/cleanup_stemcells"
	common "github.com/maximilien/bosh-softlayer-stemcells/common"
)

var _ = Describe("cleanup-stemcell command", func() {
	var (
		err                 error
		cmd                 cmds.CommandInterface
		cleanupStemcellsCmd *cleanup_stemcells.CleanupStemcellsCmd
		fakeClient          *slclientfakes.FakeSoftLayerClient
		options             common.Options
	)

	BeforeEach(func() {
		username := os.Getenv("SL_USERNAME")
		Expect(username).ToNot(Equal(""), "Missing SL_USERNAME environment variables")

		apiKey := os.Getenv("SL_API_KEY")
		Expect(apiKey).ToNot(Equal(""), "Missing SL_API_KEY environment variables")

		fakeClient = slclientfakes.NewFakeSoftLayerClient(username, apiKey)
		Expect(fakeClient).ToNot(BeNil())

		options = common.Options{
			CommandFlag:       "cleanup-stemcells",
			NamePatternFlag:   "fake-name-pattern",
			LastValidDateFlag: "2015-09-02",
			ShipItTagFlag:     "fake-shipit-tag",
		}

		cleanupStemcellsCmd = cleanup_stemcells.NewCleanupStemcellsCmd(options, fakeClient)
		cmd = cleanupStemcellsCmd
	})

	Describe("#Options", func() {
		Context("when only required options flag are passed", func() {
			BeforeEach(func() {
				options.LastValidDateFlag = ""
				options.ShipItTagFlag = ""

				cmd = cleanup_stemcells.NewCleanupStemcellsCmd(options, fakeClient)
			})

			It("returns an Options with CommandFlag and NamePatternFlag", func() {
				Expect(cmd.Options().CommandFlag).To(Equal("cleanup-stemcells"))
				Expect(cmd.Options().NamePatternFlag).To(Equal("fake-name-pattern"))
			})

			It("returns an Options with LastValidDateFlag and ShipItTagFlag as empty strings", func() {
				Expect(cmd.Options().LastValidDateFlag).To(Equal(""))
				Expect(cmd.Options().ShipItTagFlag).To(Equal(""))
			})
		})

		Context("when all options flags are passed", func() {
			It("returns an Options object with all values", func() {
				Expect(cmd.Options().CommandFlag).To(Equal("cleanup-stemcells"))
				Expect(cmd.Options().NamePatternFlag).To(Equal("fake-name-pattern"))
				Expect(cmd.Options().LastValidDateFlag).To(Equal("2015-09-02"))
				Expect(cmd.Options().ShipItTagFlag).To(Equal("fake-shipit-tag"))
			})
		})
	})

	Describe("#CheckOptions", func() {
		Context("when required options flags passed", func() {
			It("succeeds when both passed", func() {
				err := cmd.CheckOptions()
				Expect(err).ToNot(HaveOccurred())
			})

			Context("when optional flags are passed", func() {
				It("succeeds", func() {
					err := cmd.CheckOptions()
					Expect(err).ToNot(HaveOccurred())
				})
			})

			Context("when optional flags are NOT passed", func() {
				BeforeEach(func() {
					options.LastValidDateFlag = ""
					options.ShipItTagFlag = ""
					cleanupStemcellsCmd = cleanup_stemcells.NewCleanupStemcellsCmd(options, fakeClient)
					cmd = cleanupStemcellsCmd
				})

				It("succeeds", func() {
					err := cmd.CheckOptions()
					Expect(err).ToNot(HaveOccurred())
				})

				It("defaults cmd.LastValidDate to today's date", func() {
					err := cmd.CheckOptions()
					Expect(err).ToNot(HaveOccurred())
					Expect(cleanupStemcellsCmd.LastValidDate).To(Equal(time.Now().Format("2006-01-_2")))
				})

				It("defaults cmd.ShipItTag to SHIPIT", func() {
					err := cmd.CheckOptions()
					Expect(err).ToNot(HaveOccurred())
					Expect(cleanupStemcellsCmd.ShipItTag).To(Equal("SHIPIT"))
				})
			})
		})

		Context("when required options flags not passed", func() {
			It("fails when both not passed", func() {
				options.CommandFlag = ""
				options.NamePatternFlag = ""
				cmd = cleanup_stemcells.NewCleanupStemcellsCmd(options, fakeClient)

				err := cmd.CheckOptions()
				Expect(err).To(HaveOccurred())
			})

			// It("fails when CommandFlag is not cleanup-stemcells", func() {
			// 	options.CommandFlag = "fake-command"
			// 	cmd = cleanup_stemcells.NewCleanupStemcellsCmd(options, fakeClient)

			// 	err := cmd.CheckOptions()
			// 	Expect(err).To(HaveOccurred())
			// })

			// It("fails when CommandFlag not passed", func() {
			// 	options.CommandFlag = ""
			// 	cmd = cleanup_stemcells.NewCleanupStemcellsCmd(options, fakeClient)

			// 	err := cmd.CheckOptions()
			// 	Expect(err).To(HaveOccurred())
			// })

			It("fails when NamePatternFlag not passed", func() {
				options.NamePatternFlag = ""
				cmd = cleanup_stemcells.NewCleanupStemcellsCmd(options, fakeClient)

				err := cmd.CheckOptions()
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("#Run", func() {
		BeforeEach(func() {
			fakeClient.DoRawHttpRequestResponse, err = ReadJsonTestFixtures("services", "SoftLayer_Account_Service_getBlockDeviceTemplateGroups_no_objects.json")
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when no VGBDT objects are found", func() {
			It("does not issue any VGBDT delete and returns", func() {
				err := cmd.Run()
				Expect(err).ToNot(HaveOccurred())

				Expect(fakeClient.DoRawHttpRequestResponseCount).To(Equal(1))
			})
		})

		Context("when some VGBDT objects are found and they pass our filter", func() {
			BeforeEach(func() {
				vgdtgBytes, err := ReadJsonTestFixtures("services", "SoftLayer_Account_Service_getBlockDeviceTemplateGroups_some_objects.json")
				Expect(err).ToNot(HaveOccurred())

				deleteObjectBytes, err := ReadJsonTestFixtures("services", "SoftLayer_Virtual_Guest_Block_Device_Template_Group_Service_deleteObject_success.json")
				Expect(err).ToNot(HaveOccurred())

				fakeClient.DoRawHttpRequestResponses = [][]byte{vgdtgBytes, deleteObjectBytes}

				options.NamePatternFlag = "ubuntu-10.04-bosh-2400-IEM-itcs104-dea-stemcell"

				cleanupStemcellsCmd = cleanup_stemcells.NewCleanupStemcellsCmd(options, fakeClient)
				cmd = cleanupStemcellsCmd
			})

			Context("when one VGBDT object passes filter", func() {
				It("does issue one VGBDT delete and returns", func() {
					err := cmd.Run()
					Expect(err).ToNot(HaveOccurred())

					Expect(fakeClient.DoRawHttpRequestResponseCount).To(Equal(2))
				})
			})

			Context("when multiple VGBDT objects pass the filter", func() {
				BeforeEach(func() {
					vgdtgBytes, err := ReadJsonTestFixtures("services", "SoftLayer_Account_Service_getBlockDeviceTemplateGroups_some_objects.json")
					Expect(err).ToNot(HaveOccurred())

					deleteObjectBytes, err := ReadJsonTestFixtures("services", "SoftLayer_Virtual_Guest_Block_Device_Template_Group_Service_deleteObject_success.json")
					Expect(err).ToNot(HaveOccurred())

					fakeClient.DoRawHttpRequestResponses = [][]byte{vgdtgBytes, deleteObjectBytes, deleteObjectBytes}

					options = common.Options{
						CommandFlag:       "cleanup-stemcells",
						LastValidDateFlag: "2015-09-01",
						NamePatternFlag:   "ubuntu-10.04-bosh-[0-9]*-IEM-itcs104-dea-stemcell",
					}

					cleanupStemcellsCmd = cleanup_stemcells.NewCleanupStemcellsCmd(options, fakeClient)
					cmd = cleanupStemcellsCmd
				})
				It("does issue 2 VGBDT deltes and returns", func() {

					err := cmd.Run()
					Expect(err).ToNot(HaveOccurred())

					Expect(fakeClient.DoRawHttpRequestResponseCount).To(Equal(3))
				})
			})
		})
	})

	Describe("#FilterVGBDTG", func() {
		BeforeEach(func() {
			vgdtgBytes, err := ReadJsonTestFixtures("services", "SoftLayer_Account_Service_getBlockDeviceTemplateGroups_some_objects.json")
			Expect(err).ToNot(HaveOccurred())

			fakeClient.DoRawHttpRequestResponse = vgdtgBytes
		})

		Describe("#FilterVGBDTGTag", func() {
			Context("when no VGBDT objects are filtered out", func() {
				It("filters out NONE of the VGBDT objects", func() {
					accountService, err := fakeClient.GetSoftLayer_Account_Service()
					Expect(err).ToNot(HaveOccurred())

					vgbdtgObjects, err := accountService.GetBlockDeviceTemplateGroups()
					Expect(err).ToNot(HaveOccurred())

					filteredObjects := cleanup_stemcells.FilterVGBDTGTag(vgbdtgObjects, "NotSHIPIT")
					Expect(err).ToNot(HaveOccurred())

					Expect(len(filteredObjects)).To(Equal(4))
				})
			})

			Context("when some VGBDT objects are filtered out", func() {
				It("filters out two objects that have the SHIPIT tag", func() {
					accountService, err := fakeClient.GetSoftLayer_Account_Service()
					Expect(err).ToNot(HaveOccurred())

					vgbdtgObjects, err := accountService.GetBlockDeviceTemplateGroups()
					Expect(err).ToNot(HaveOccurred())

					filteredObjects := cleanup_stemcells.FilterVGBDTGTag(vgbdtgObjects, "SHIPIT")
					Expect(err).ToNot(HaveOccurred())

					Expect(len(filteredObjects)).To(Equal(2))

					foundObject := false
					for _, object := range filteredObjects {
						if object.Id == 2 {
							foundObject = true
						}
						// fmt.Sprintf("=======>object.Id id:%d", object.Id)
					}
					Expect(foundObject).To(Equal(true))

					foundObject = false
					for _, object := range filteredObjects {
						if object.Id == 3 {
							foundObject = true
						}
					}
					Expect(foundObject).To(Equal(true))

				})
			})
		})

		Describe("#FilterVGBDTGDate", func() {
			Context("when no VGBDT objects are filtered out", func() {
				It("filters NONE of the VGBDT objects", func() {
					accountService, err := fakeClient.GetSoftLayer_Account_Service()
					Expect(err).ToNot(HaveOccurred())

					vgbdtgObjects, err := accountService.GetBlockDeviceTemplateGroups()
					Expect(err).ToNot(HaveOccurred())

					filteredObjects, err := cleanup_stemcells.FilterVGBDTGDate(vgbdtgObjects, time.Now().Format("2006/01/_2"))
					Expect(err).ToNot(HaveOccurred())

					Expect(len(filteredObjects)).To(Equal(4))
				})
			})

			Context("when some VGBDT objects are filtered out", func() {
				It("succeeds when 2 objects are filtered out that were created after July 10th, 2015", func() {
					accountService, err := fakeClient.GetSoftLayer_Account_Service()
					Expect(err).ToNot(HaveOccurred())

					vgbdtgObjects, err := accountService.GetBlockDeviceTemplateGroups()
					Expect(err).ToNot(HaveOccurred())

					filteredObjects, err := cleanup_stemcells.FilterVGBDTGDate(vgbdtgObjects, "2015-07-10")
					Expect(err).ToNot(HaveOccurred())

					Expect(len(filteredObjects)).To(Equal(2))

					foundObject := false
					for _, object := range filteredObjects {
						if object.Id == 1 {
							foundObject = true
						}
						// fmt.Sprintf("=======>object.Id id:%d", object.Id)
					}
					Expect(foundObject).To(Equal(true))

					foundObject = false
					for _, object := range filteredObjects {
						if object.Id == 2 {
							foundObject = true
						}
					}
					Expect(foundObject).To(Equal(true))
				})

				It("succeeds when 3 objects passes filter that as created after August 25th, 2015", func() {
					accountService, err := fakeClient.GetSoftLayer_Account_Service()
					Expect(err).ToNot(HaveOccurred())

					vgbdtgObjects, err := accountService.GetBlockDeviceTemplateGroups()
					Expect(err).ToNot(HaveOccurred())

					filteredObjects, err := cleanup_stemcells.FilterVGBDTGDate(vgbdtgObjects, "2015-08-25")
					Expect(err).ToNot(HaveOccurred())

					Expect(len(filteredObjects)).To(Equal(3))
					foundObject := false
					for _, object := range filteredObjects {
						if object.Id == 1 {
							foundObject = true
						}
						// fmt.Sprintf("=======>object.Id id:%d", object.Id)
					}
					Expect(foundObject).To(Equal(true))

					foundObject = false
					for _, object := range filteredObjects {
						if object.Id == 2 {
							foundObject = true
						}
					}
					Expect(foundObject).To(Equal(true))

					foundObject = false
					for _, object := range filteredObjects {
						if object.Id == 3 {
							foundObject = true
						}
					}
					Expect(foundObject).To(Equal(true))
				})

			})

			Context("when ALL VGBDT objects are made after given date", func() {
				It("filters out all VGBDT objects", func() {
					accountService, err := fakeClient.GetSoftLayer_Account_Service()
					Expect(err).ToNot(HaveOccurred())

					vgbdtgObjects, err := accountService.GetBlockDeviceTemplateGroups()
					Expect(err).ToNot(HaveOccurred())

					filteredObjects, err := cleanup_stemcells.FilterVGBDTGDate(vgbdtgObjects, "2010-08-25")
					Expect(err).ToNot(HaveOccurred())

					Expect(len(filteredObjects)).To(Equal(0))
				})
			})

			Context("when the date string is not formatted correctly", func() {
				It("fails when the date string given is not in 'year-month-day' format", func() {
					accountService, err := fakeClient.GetSoftLayer_Account_Service()
					Expect(err).ToNot(HaveOccurred())

					vgbdtgObjects, err := accountService.GetBlockDeviceTemplateGroups()
					Expect(err).ToNot(HaveOccurred())

					_, err = cleanup_stemcells.FilterVGBDTGDate(vgbdtgObjects, "01-08-2015")
					Expect(err).To(HaveOccurred())
				})

				It("fails when the date string given is not in 'year-month-day' or 'year/month/day' format", func() {
					accountService, err := fakeClient.GetSoftLayer_Account_Service()
					Expect(err).ToNot(HaveOccurred())

					vgbdtgObjects, err := accountService.GetBlockDeviceTemplateGroups()
					Expect(err).ToNot(HaveOccurred())

					_, err = cleanup_stemcells.FilterVGBDTGDate(vgbdtgObjects, "01/08-2015")
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Describe("#FilterVGBDTGName", func() {
			Context("when some VGBDT objects match the pattern", func() {
				It("filters three VGBT objects matching name: ubuntu-10.04-bosh-[0-9][0-9][0-9]+-IEM-itcs104-dea-stemcell", func() {
					accountService, err := fakeClient.GetSoftLayer_Account_Service()
					Expect(err).ToNot(HaveOccurred())

					vgbdtgObjects, err := accountService.GetBlockDeviceTemplateGroups()
					Expect(err).ToNot(HaveOccurred())

					filteredObjects, err := cleanup_stemcells.FilterVGBDTGName(vgbdtgObjects, "ubuntu-10.04-bosh-[0-9][0-9][0-9]+-IEM-itcs104-dea-stemcell")
					Expect(err).ToNot(HaveOccurred())

					Expect(len(filteredObjects)).To(Equal(3))

					foundObject := false
					for _, object := range filteredObjects {
						if object.Id == 1 {
							foundObject = true
						}
						// fmt.Sprintf("=======>object.Id id:%d", object.Id)
					}
					Expect(foundObject).To(Equal(true))

					foundObject = false
					for _, object := range filteredObjects {
						if object.Id == 2 {
							foundObject = true
						}
					}
					Expect(foundObject).To(Equal(true))

					foundObject = false
					for _, object := range filteredObjects {
						if object.Id == 3 {
							foundObject = true
						}
					}
					Expect(foundObject).To(Equal(true))
				})
			})

			Context("when ALL VGBDT objects match the pattern", func() {
				It("filters 0 VGBDT objects matching name: ubuntu-10.04-bosh-0-IEM-itcs104-dea-stemcell", func() {
					accountService, err := fakeClient.GetSoftLayer_Account_Service()
					Expect(err).ToNot(HaveOccurred())

					vgbdtgObjects, err := accountService.GetBlockDeviceTemplateGroups()
					Expect(err).ToNot(HaveOccurred())

					filteredObjects, err := cleanup_stemcells.FilterVGBDTGName(vgbdtgObjects, "ubuntu-10.04-bosh-0-IEM-itcs104-dea-stemcell")
					Expect(err).ToNot(HaveOccurred())

					Expect(len(filteredObjects)).To(Equal(0))
				})
			})

			Context("when the name pattern is not of regex format", func() {
				It("fails when the name pattern does not follow a regex pattern", func() {
					accountService, err := fakeClient.GetSoftLayer_Account_Service()
					Expect(err).ToNot(HaveOccurred())

					vgbdtgObjects, err := accountService.GetBlockDeviceTemplateGroups()
					Expect(err).ToNot(HaveOccurred())

					_, err = cleanup_stemcells.FilterVGBDTGDate(vgbdtgObjects, "ubuntu-10.04-bosh-0*+-IEM-itcs104-dea-stemcell")
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("when no objects pass the combined filter", func() {
			BeforeEach(func() {
				options.LastValidDateFlag = "2015-07-20"
				options.ShipItTagFlag = "fake-tag"
				options.NamePatternFlag = "ubuntu-10.04-bosh-[0-9][0-9]-IEM-itcs104-dea-stemcell"

				cleanupStemcellsCmd = cleanup_stemcells.NewCleanupStemcellsCmd(options, fakeClient)
			})

			It("filters out ALL objects out with the combination of the filters ", func() {
				accountService, err := fakeClient.GetSoftLayer_Account_Service()
				Expect(err).ToNot(HaveOccurred())

				vgbdtgObjects, err := accountService.GetBlockDeviceTemplateGroups()
				Expect(err).ToNot(HaveOccurred())

				filteredObjects, err := cleanupStemcellsCmd.FilterVGBDTG(vgbdtgObjects)
				Expect(err).ToNot(HaveOccurred())

				Expect(len(filteredObjects)).To(Equal(0))
			})

		})

		Context("when all objects pass the combined filter", func() {
			It("filters NO objects out with the combination of the filters", func() {

				options.LastValidDateFlag = "2015-09-01"
				options.ShipItTagFlag = "NotSHIPIT"
				options.NamePatternFlag = "ubuntu-10.04-bosh-[0-9]+-IEM-itcs104-dea-stemcell"

				cleanupStemcellsCmd = cleanup_stemcells.NewCleanupStemcellsCmd(options, fakeClient)

				accountService, err := fakeClient.GetSoftLayer_Account_Service()
				Expect(err).ToNot(HaveOccurred())

				vgbdtgObjects, err := accountService.GetBlockDeviceTemplateGroups()
				Expect(err).ToNot(HaveOccurred())

				filteredObjects, err := cleanupStemcellsCmd.FilterVGBDTG(vgbdtgObjects)
				Expect(err).ToNot(HaveOccurred())

				Expect(len(filteredObjects)).To(Equal(4))
			})

			It("filters NONE of the objects due to default: LastValidDateFlag = Today ", func() {
				options.LastValidDateFlag = ""
				options.ShipItTagFlag = "NotSHIPIT"
				options.NamePatternFlag = "ubuntu-10.04-bosh-[0-9]*-IEM-itcs104-dea-stemcell"

				cleanupStemcellsCmd = cleanup_stemcells.NewCleanupStemcellsCmd(options, fakeClient)

				accountService, err := fakeClient.GetSoftLayer_Account_Service()
				Expect(err).ToNot(HaveOccurred())

				vgbdtgObjects, err := accountService.GetBlockDeviceTemplateGroups()
				Expect(err).ToNot(HaveOccurred())

				filteredObjects, err := cleanupStemcellsCmd.FilterVGBDTG(vgbdtgObjects)
				Expect(err).ToNot(HaveOccurred())

				Expect(len(filteredObjects)).To(Equal(4))
			})
		})

		Context("when some objects pass the combined filter", func() {
			It("filters two objects out with default SHIPIT tag: SHIPIT", func() {
				options = common.Options{
					CommandFlag:       "cleanup-stemcells",
					LastValidDateFlag: "2015-09-01",
					NamePatternFlag:   "ubuntu-10.04-bosh-[0-9]*-IEM-itcs104-dea-stemcell",
				}

				cleanupStemcellsCmd = cleanup_stemcells.NewCleanupStemcellsCmd(options, fakeClient)

				accountService, err := fakeClient.GetSoftLayer_Account_Service()
				Expect(err).ToNot(HaveOccurred())

				vgbdtgObjects, err := accountService.GetBlockDeviceTemplateGroups()
				Expect(err).ToNot(HaveOccurred())

				filteredObjects, err := cleanupStemcellsCmd.FilterVGBDTG(vgbdtgObjects)
				Expect(err).ToNot(HaveOccurred())

				Expect(len(filteredObjects)).To(Equal(2))

				foundObject := false
				for _, object := range filteredObjects {
					if object.Id == 2 {
						foundObject = true
					}
					// fmt.Sprintf("=======>object.Id id:%d", object.Id)
				}
				Expect(foundObject).To(Equal(true))

				foundObject = false
				for _, object := range filteredObjects {
					if object.Id == 3 {
						foundObject = true
					}
				}
				Expect(foundObject).To(Equal(true))
			})

			It("filters out 3 of the objects out with the combined filter", func() {
				options = common.Options{
					CommandFlag:       "cleanup-stemcells",
					LastValidDateFlag: "2015-08-30",
					NamePatternFlag:   "ubuntu-10.04-bosh-23[0-9]*-IEM-itcs104-dea-stemcell",
					ShipItTagFlag:     "SHIPIT",
				}

				cleanupStemcellsCmd = cleanup_stemcells.NewCleanupStemcellsCmd(options, fakeClient)

				accountService, err := fakeClient.GetSoftLayer_Account_Service()
				Expect(err).ToNot(HaveOccurred())

				vgbdtgObjects, err := accountService.GetBlockDeviceTemplateGroups()
				Expect(err).ToNot(HaveOccurred())

				filteredObjects, err := cleanupStemcellsCmd.FilterVGBDTG(vgbdtgObjects)
				Expect(err).ToNot(HaveOccurred())

				Expect(len(filteredObjects)).To(Equal(1))
				Expect(filteredObjects[0].Id).To(Equal(3))

			})
		})

	})
})
