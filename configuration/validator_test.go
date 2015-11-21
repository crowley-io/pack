package configuration

// func TestValidateConfigurationEmptyOutput(t *testing.T) {
//
// 	c := &configuration.Configuration{
// 		Install: configuration.Install{
// 			Path:  "/root",
// 			Image: "debian",
// 		},
// 	}
// 	err := ValidateConfiguration(c)
//
// 	assert.NotNil(t, err)
// 	assert.Equal(t, ErrOutputRequired, err)
//
// }
//
// func TestValidateConfigurationEmptyImage(t *testing.T) {
//
// 	c := &configuration.Configuration{
// 		Output: "file",
// 		Install: configuration.Install{
// 			Path: "/root",
// 		},
// 	}
// 	err := ValidateConfiguration(c)
//
// 	assert.NotNil(t, err)
// 	assert.Equal(t, ErrImageRequired, err)
//
// }
//
// func TestValidateConfigurationEmptyPath(t *testing.T) {
//
// 	c := &configuration.Configuration{
// 		Output: "file",
// 		Install: configuration.Install{
// 			Image: "debian",
// 		},
// 	}
// 	err := ValidateConfiguration(c)
//
// 	assert.NotNil(t, err)
// 	assert.Equal(t, ErrPathRequired, err)
//
// }
