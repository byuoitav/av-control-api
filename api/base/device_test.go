package base_test

import (
	"encoding/json"
	"testing"

	"github.com/byuoitav/av-control-api/api/base"
	"github.com/matryer/is"
)

var validDeviceStruct = base.Device{
	ID:      "ITB-1101-SW1",
	Name:    "ITB-1101-SW1",
	Address: "127.0.0.1",
	Type: base.DeviceType{
		ID:          "Atlona4x1",
		Description: "Atlona 4x1 Video Switcher",
		DisplayName: "Atlona Video Switcher - 4x1",
		Source:      true,
		Destination: true,
		Roles:       []base.Role{},
		Ports:       []base.Port{},
		Commands: map[string]base.Command{
			"ChangeInput": base.Command{
				Addresses: map[string]string{
					"fallback": "http://localhost:8026/{{.Address}}/output/{{.Output}}/input/{{.Input}}",
					"aws":      "http://atlona-driver-prd/{{.Address}}/output/{{.Output}}/input/{{.Input}}",
				},
				Order: 10,
			},
			"InputStatus": {
				Addresses: map[string]string{
					"fallback": "http://localhost:8026/{{.Address}}/output/{{.Port}}",
					"aws":      "http://atlona-driver-prd/{{.Address}}/output/{{.Port}}",
				},
				Order: 10,
			},
			"NoFallback": {
				Addresses: map[string]string{
					"aws": "http://atlona-driver-prd/{{.Address}}/output/{{.Port}}",
				},
				Order: 10,
			},
		},
		Tags: []string{},
	},
	Proxy: map[string]string{
		"InputStatus": "proxyhost:8080",
	},
}

var validDeviceTypeDocument = `
{
  "_id": "Atlona4x1",
  "description": "Atlona 4x1 Video Switcher",
  "display_name": "Atlona Video Switcher - 4x1",
  "source": true,
  "destination": true,
  "roles": [],
  "ports": [],
  "commands": {
    "ChangeInput": {
        "addresses": {
          "fallback": "http://localhost:8026/{{.Address}}/output/{{.Output}}/input/{{.Input}}",
          "aws": "http://atlona-driver-prd/{{.Address}}/output/{{.Output}}/input/{{.Input}}"
        },
        "order": 10
    },
    "InputStatus": {
        "addresses": {
          "fallback": "http://localhost:8026/{{.Address}}/output/{{.Port}}",
          "aws": "http://atlona-driver-prd/{{.Address}}/output/{{.Port}}"
        },
        "order": 10
    },
    "NoFallback": {
        "addresses": {
          "aws": "http://atlona-driver-prd/{{.Address}}/output/{{.Port}}"
        },
        "order": 10
    }
  },
  "tags": []
}
`

var validCommandDocument = `
{
	"addresses": {
		"fallback": "http://localhost:8026/{{.Address}}/output/{{.Output}}/input/{{.Input}}",
		"aws": "http://atlona-driver-prd/{{.Address}}/output/{{.Output}}/input/{{.Input}}"
	},
	"order": 10
}
`

func TestCommandUnmarshal(t *testing.T) {
	is := is.New(t)

	t.Run("Addresses Unmarshal Correctly", func(t *testing.T) {
		is := is.New(t)
		c := base.Command{}
		err := json.Unmarshal([]byte(validCommandDocument), &c)
		is.NoErr(err)
		is.Equal(c, validDeviceStruct.Type.Commands["ChangeInput"]) // The unmarshalled struct should match the expected struct
	})
}

func TestDeviceTypeUnmarshal(t *testing.T) {
	is := is.New(t)

	t.Run("Commands unmarshal successfully", func(t *testing.T) {
		is := is.New(t)
		dt := base.DeviceType{}
		err := json.Unmarshal([]byte(validDeviceTypeDocument), &dt)
		is.NoErr(err)
		is.Equal(dt, validDeviceStruct.Type) // The unmarshalled struct should match the expected struct

	})
}

func TestGetCommandByID(t *testing.T) {
	is := is.New(t)

	t.Run("Should return command for valid id", func(t *testing.T) {
		is := is.New(t)
		c, err := validDeviceStruct.GetCommandByID("ChangeInput")
		is.NoErr(err)                                               // Expected to run without error
		is.Equal(c, validDeviceStruct.Type.Commands["ChangeInput"]) // Expected to get the valid command struct
	})

	t.Run("Should return an empty command and an error for an invalid id", func(t *testing.T) {
		is := is.New(t)
		expected := base.Command{}
		got, err := validDeviceStruct.GetCommandByID("foo")
		is.True(err != nil)     // Expected to get an error back
		is.Equal(got, expected) // Should have returned an empty command on an invalid command ID
	})
}

func TestHasCommand(t *testing.T) {
	is := is.New(t)

	t.Run("Should return true for a valid command ID", func(t *testing.T) {
		is := is.New(t)
		is.True(validDeviceStruct.HasCommand("ChangeInput"))
	})

	t.Run("Should return false for an invalid command ID", func(t *testing.T) {
		is := is.New(t)
		is.True(!validDeviceStruct.HasCommand("foo"))
	})
}

func TestCommandValidate(t *testing.T) {
	is := is.New(t)

	t.Run("Should successfully validate a valid command", func(t *testing.T) {
		is := is.New(t)
		err := validDeviceStruct.Type.Commands["ChangeInput"].Validate()
		is.NoErr(err) // Expected struct to validate successfully
	})

	t.Run("Should fail to validate a command with malformed address", func(t *testing.T) {
		is := is.New(t)
		err := base.Command{
			Addresses: map[string]string{
				"fallback": "http://local^host:8020",
			},
		}.Validate()
		is.True(err != nil)
	})
}

func TestDeviceTypeValidate(t *testing.T) {
	is := is.New(t)

	t.Run("Should successfully validate a valid Device Type", func(t *testing.T) {
		is := is.New(t)
		err := validDeviceStruct.Type.Validate(true)
		is.NoErr(err) // Expected struct to validate successfully
	})

	t.Run(
		"Should fail to validate if a Command ID is less than three characters",
		func(t *testing.T) {
			is := is.New(t)
			badType := base.DeviceType{
				ID: "BadDevice",
				Commands: map[string]base.Command{
					"AB": base.Command{},
				},
			}
			err := badType.Validate(true)
			is.True(err != nil)
		},
	)
}

func TestBuildCommandURL(t *testing.T) {
	is := is.New(t)

	t.Run(
		"Should return the requested url when the env exists as an address and no proxy",
		func(t *testing.T) {
			is := is.New(t)
			expected := validDeviceStruct.Type.Commands["ChangeInput"].Addresses["aws"]
			got, err := validDeviceStruct.BuildCommandURL("ChangeInput", "aws")
			is.NoErr(err)           // Expected to run without error
			is.Equal(got, expected) // Expected the returned string to match the expected
		},
	)
	t.Run(
		"Should return the fallback url when the env doesn't exist as an address and no proxy",
		func(t *testing.T) {
			is := is.New(t)
			expected := validDeviceStruct.Type.Commands["ChangeInput"].Addresses["fallback"]
			got, err := validDeviceStruct.BuildCommandURL("ChangeInput", "foobar")
			is.NoErr(err)           // Expected to run without error
			is.Equal(got, expected) // Expected the returned string to match the expected
		},
	)
	t.Run(
		"Should return the proxied url when the env exists and there is a proxy",
		func(t *testing.T) {
			is := is.New(t)
			expected := "http://proxyhost:8080/{{.Address}}/output/{{.Port}}"
			got, err := validDeviceStruct.BuildCommandURL("InputStatus", "aws")
			is.NoErr(err)           // Expected to run without error
			is.Equal(got, expected) // Expected the returned string to match the expected
		},
	)
	t.Run(
		"Should return the proxied url when the env doesn't exist (but there is a fallback) and there is a proxy",
		func(t *testing.T) {
			is := is.New(t)
			expected := "http://proxyhost:8080/{{.Address}}/output/{{.Port}}"
			got, err := validDeviceStruct.BuildCommandURL("InputStatus", "foobar")
			is.NoErr(err)           // Expected to run without error
			is.Equal(got, expected) // Expected the returned string to match the expected
		},
	)
	t.Run(
		"Should return the requested address even if no fallback is present",
		func(t *testing.T) {
			is := is.New(t)
			expected := validDeviceStruct.Type.Commands["NoFallback"].Addresses["aws"]
			got, err := validDeviceStruct.BuildCommandURL("NoFallback", "aws")
			is.NoErr(err)           // Expected to run without error
			is.Equal(got, expected) // Expected the returned string to match the expected
		},
	)
	t.Run(
		"Should return an error if the env doesn't exist and there is no fallback address",
		func(t *testing.T) {
			is := is.New(t)
			_, err := validDeviceStruct.BuildCommandURL("NoFallback", "foobar")
			is.True(err != nil) // Expected to get an error
		},
	)

}
