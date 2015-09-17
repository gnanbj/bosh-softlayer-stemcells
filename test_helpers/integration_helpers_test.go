package test_helpers_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	test_helpers "github.com/maximilien/bosh-softlayer-stemcells/test_helpers"
)

var _ = Describe("Helper Functions for running Integration Tests", func() {
	Context("#RunSLStemcells", func() {
		It("runs /out/sl_stemcells with help option and has correct output", func() {
			pwd, err := os.Getwd()
			Expect(err).ToNot(HaveOccurred())

			rootPath := filepath.Join(pwd, "..")

			output, err := test_helpers.RunSLStemcells(rootPath, []string{"--help"})
			Expect(err).ToNot(HaveOccurred())

			Expect(string(output)).To(ContainSubstring(`
usage: bosh-softlayer-stemcells -c <command> [--name <template-name>] [--note <import note>] 
       --os-ref-code <OsRefCode> --uri <swiftURI>

  -h | --help   prints the usage

  IMPORT-IMAGE:

  -c import-image  the import image command
  --name           the group name to be applied to the imported template
  --note           the note to be applied to the imported template
  --os-ref-code    the referenceCode of the operating system software 
                   description for the imported VHD 
                   available options: CENTOS_6_32, CENTOS_6_64, CENTOS_7_64, 
                     REDHAT_6_32, REDHAT_6_64, REDHAT_7_64, UBUNTU_10_32, 
                     UBUNTU_10_64, UBUNTU_12_32, UBUNTU_12_64, UBUNTU_14_32, 
                     UBUNTU_14_64, WIN_2003-STD-SP2-5_32, WIN_2003-STD-SP2-5_64, 
                     WIN_2012-STD_64
  --uri            the URI for an object storage object (.vhd/.iso file)
                   swift://<ObjectStorageAccountName>@<clusterName>/<containerName>/<fileName.(vhd|iso)>`))
		})

		Context("#RunSLStemcells with cleanup-stemcells command", func() {
			It("runs /out/sl_stemcells cleanup-stemcells with namePattern and date", func() {
				pwd, err := os.Getwd()
				Expect(err).ToNot(HaveOccurred())

				rootPath := filepath.Join(pwd, "..")

				output, err := test_helpers.RunSLStemcells(rootPath, []string{"-c", "cleanup-stemcells", "--name-pattern", "fake-name"})
				Expect(err).ToNot(HaveOccurred())

				Expect(string(output)).To(ContainSubstring("Total time:"))

			})
		})
	})

})
