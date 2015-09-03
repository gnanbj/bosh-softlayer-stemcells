package cleanup_stemcells_test

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	slclientfakes "github.com/maximilien/softlayer-go/client/fakes"

	cmds "github.com/maximilien/bosh-softlayer-stemcells/cmds"
	cleanup_stemcells "github.com/maximilien/bosh-softlayer-stemcells/cmds/cleanup_stemcells"
	common "github.com/maximilien/bosh-softlayer-stemcells/common"
)

var _ = Describe("cleanup-stemcell command", func() {
	var (
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
					Expect(cleanupStemcellsCmd.LastValidDate).To(Equal(time.Now().Format(time.UnixDate)))
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

			It("fails CommandFlag is not cleanup-stemcells", func() {
				options.CommandFlag = "fake-command"
				cmd = cleanup_stemcells.NewCleanupStemcellsCmd(options, fakeClient)

				err := cmd.CheckOptions()
				Expect(err).To(HaveOccurred())
			})

			It("fails CommandFlag not passed", func() {
				options.CommandFlag = ""
				cmd = cleanup_stemcells.NewCleanupStemcellsCmd(options, fakeClient)

				err := cmd.CheckOptions()
				Expect(err).To(HaveOccurred())
			})

			It("fails NamePatternFlag not passed", func() {
				options.NamePatternFlag = ""
				cmd = cleanup_stemcells.NewCleanupStemcellsCmd(options, fakeClient)

				err := cmd.CheckOptions()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
