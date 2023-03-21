resource "synthetics_create_variable_v2" "variable_v2_foo" {
  variable {
    description = "The most awesome variable. Full of snakes."
    value = "barv3--oopsasdasd"
    // Once created name and secret can not be changed and will result in a 422 from the API
    // unless the variable is deleted and re-created
    name = "terraform-test121"
    secret = false  
  }    
}