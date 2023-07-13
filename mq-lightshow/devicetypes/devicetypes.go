// Package devicetypes provides our device types.
package devicetypes

import (
	"github.com/lovesway/hassio-addons/mq-lightshow/models"
)

// DeviceTypes represents the controller for the API.
type DeviceTypes struct {
	deviceTypes []models.DeviceType
}

// NewDeviceTypes provides an instance of DeviceTypes.
func NewDeviceTypes() DeviceTypes {
	dts := []models.DeviceType{}

	// Manually set up our device types which will seldom change.
	dt := models.DeviceType{
		ID:       1,
		Name:     "Tasmota",
		Commands: getTasmotaCommands(),
	}

	dts = append(dts, dt)

	return DeviceTypes{
		deviceTypes: dts,
	}
}

// GetDeviceTypes will return our deviceTypes.
func (d DeviceTypes) GetDeviceTypes() []models.DeviceType {
	return d.deviceTypes
}

func getTasmotaCommands() []models.Command {
	commands := []models.Command{}
	c := models.Command{
		Name:        "Power",
		Description: "Power - Toggle, On, Off",
	}
	commands = append(commands, c)

	c = models.Command{
		Name:        "Color1",
		Description: "Color1 - Set color, values can be r,g,b or #hex",
	}
	commands = append(commands, c)

	c = models.Command{
		Name:        "Color2",
		Description: "Color2 - Set color adjusted to current Dimmer value",
	}
	commands = append(commands, c)

	c = models.Command{
		Name:        "CT",
		Description: "CT - Set color temperature from 153 (cold) to 500 (warm) for CT lights",
	}
	commands = append(commands, c)

	c = models.Command{
		Name:        "Dimmer",
		Description: "Dimmer - Set dimmer value from 0 to 100(%)",
	}
	commands = append(commands, c)

	c = models.Command{
		Name:        "Fade",
		Description: "Fade - 0 = do not use fade (default), 1 = use fade",
	}
	commands = append(commands, c)

	c = models.Command{
		Name:        "Speed",
		Description: "Speed - Set fade speed from fast 1 to very slow 40 (The Speed value represents the time in 0.5s)",
	}
	commands = append(commands, c)

	c = models.Command{
		Name:        "HsbColor",
		Description: "HsbColor - hue,sat,bri = set color by hue, saturation and brightness",
	}
	commands = append(commands, c)

	c = models.Command{
		Name:        "HsbColor1",
		Description: "HsbColor1 - 0..360 = set hue",
	}
	commands = append(commands, c)

	c = models.Command{
		Name:        "HsbColor2",
		Description: "HsbColor2 - 0..100 = set saturation",
	}
	commands = append(commands, c)

	c = models.Command{
		Name:        "HsbColor3",
		Description: "HsbColor3 - 0..100 = set brightness",
	}
	commands = append(commands, c)

	c = models.Command{
		Name: "Scheme",
		Description: "Scheme - 0 = single color, 1 = start wake up, 2 = cycle up colors, " +
			"3 = cycle down colors, 4 = random cycle colors",
	}
	commands = append(commands, c)

	c = models.Command{
		Name: "Wakeup",
		Description: "Wakeup - Start wake up from OFF to stored Dimmer value " +
			"(0..100 = Start wake up from OFF to provided value)",
	}
	commands = append(commands, c)

	c = models.Command{
		Name:        "WakeupDuration",
		Description: "WakeupDuration - 1..3000 = set wake up duration in seconds",
	}
	commands = append(commands, c)

	c = models.Command{
		Name:        "White",
		Description: "White - 1..100 = set white channel brightness in single white channel lights (single W or RGBW lights)",
	}
	commands = append(commands, c)

	return commands
}
