package main


import (
  "flag"
)


type Config struct {
  JsonFile string
}


func ParseFlags() (Config,error) {

  var config = Config{}

  var flagSet = flag.NewFlagSet("select", flag.ExitOnError )

  var jsonFile = flagSet.String("json","","JSON file to import")


  err := flagSet.Parse( flag.Args() )
  if err != nil {
    return config, err
  }

  config.JsonFile = *jsonFile


  return config, nil
}



func main() {

}
