package cleanup_stemcells

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	sldatatypes "github.com/maximilien/softlayer-go/data_types"
	softlayer "github.com/maximilien/softlayer-go/softlayer"

	common "github.com/maximilien/bosh-softlayer-stemcells/common"
)

type CleanupStemcellsCmd struct {
	LastValidDate string
	ShipItTag     string

	client  softlayer.Client
	options common.Options
}

func NewCleanupStemcellsCmd(options common.Options, client softlayer.Client) *CleanupStemcellsCmd {
	return &CleanupStemcellsCmd{
		options: options,
		client:  client,
	}
}

func (cmd *CleanupStemcellsCmd) Println(a ...interface{}) (int, error) {
	return fmt.Println(a)
}

func (cmd *CleanupStemcellsCmd) Printf(msg string, a ...interface{}) (int, error) {
	return fmt.Printf(msg, a)
}

func (cmd *CleanupStemcellsCmd) Options() common.Options {
	return cmd.options
}

func (cmd *CleanupStemcellsCmd) CheckOptions() error {
	if cmd.options.CommandFlag != "cleanup-stemcells" {
		return errors.New(fmt.Sprintf("CommandFlag (%s) is incorrect, is not cleanup-stemcells", cmd.options.CommandFlag))
	}

	if cmd.options.NamePatternFlag == "" {
		return errors.New("NamePatternFlag not passed")
	}

	cmd.defaultOptionalFlags()

	return nil
}

func (cmd *CleanupStemcellsCmd) Run() error {
	accountService, err := cmd.client.GetSoftLayer_Account_Service()
	if err != nil {
		return errors.New(fmt.Sprintf("Getting softlayer account service, message: %s", err.Error()))
	}

	vgbdtgObjects, err := accountService.GetBlockDeviceTemplateGroups()
	if err != nil {
		return errors.New(fmt.Sprintf("Getting Block Device Template Groups, message: %s", err.Error()))
	}

	toBeDeletedObjects, err := cmd.FilterVGBDTG(vgbdtgObjects)
	if err != nil {
		return errors.New(fmt.Sprintf("Filtering softlayer VGBDTG objects, message: %s", err.Error()))
	}

	vgbdtgService, err := cmd.client.GetSoftLayer_Virtual_Guest_Block_Device_Template_Group_Service()
	if err != nil {
		return errors.New(fmt.Sprintf("Getting softlayer VGBDTG service, message: %s", err.Error()))
	}

	for object, _ := range toBeDeletedObjects {
		_, err := vgbdtgService.DeleteObject(object)
		if err != nil {
			return errors.New(fmt.Sprintf("Deleting VGBDTG object, message: %s", err.Error()))
		}

		//Wait for transaction...
	}

	return nil
}

// Utility methods

func (cmd *CleanupStemcellsCmd) FilterVGBDTG(vgbdtgObjects []sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group) ([]sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group, error) {
	//filter
	// 1. by shipit tag
	// 2. by date
	// 3. by name

	return vgbdtgObjects, nil
}

func FilterVGBDTGTag(vgbdtgObjects []sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group, tag string) []sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group {
	filteredObjects := []sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group{}

	for _, object := range vgbdtgObjects {
		if object.Note != tag {
			filteredObjects = append(filteredObjects, object)
		}
	}

	return filteredObjects
}

func FilterVGBDTGDate(vgbdtgObjects []sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group, date string) ([]sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group, error) {
	filteredObjects := []sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group{}

	filterDate, err := time.Parse(time.UnixDate, date)
	if err != nil {
		return filteredObjects, errors.New(fmt.Sprintf("Converting Date into time object, message: %s", err.Error()))
	}

	for _, object := range vgbdtgObjects {
		if filterDate.After(*object.CreateDate) {
			filteredObjects = append(filteredObjects, object)
		}

	}

	return filteredObjects, nil
}

func FilterVGBDTGName(vgbdtgObjects []sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group, namePattern string) ([]sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group, error) {
	filteredObjects := []sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group{}

	compiledPattern, err := regexp.Compile(namePattern)
	if err != nil {
		return filteredObjects, errors.New(fmt.Sprintf("Compiling name pattern, message: %s", err.Error()))
	}

	for _, object := range vgbdtgObjects {
		if compiledPattern.MatchString(object.Name) {
			filteredObjects = append(filteredObjects, object)
		}
	}

	return filteredObjects, nil

}

// Private

func (cmd *CleanupStemcellsCmd) defaultOptionalFlags() {
	if cmd.options.LastValidDateFlag == "" {
		cmd.LastValidDate = time.Now().Format(time.UnixDate)
	}

	if cmd.options.ShipItTagFlag == "" || cmd.options.ShipItTagFlag == "SHIPIT" {
		cmd.ShipItTag = "SHIPIT"
	}
}
