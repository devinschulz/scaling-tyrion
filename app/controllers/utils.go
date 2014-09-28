package controllers

import "github.com/revel/revel"

func checkINFO(err error) {
  if err != nil {
    revel.INFO.Println(err)
  }
}

func checkWARN(err error) {
  if err != nil {
    revel.WARN.Println(err)
  }
}

func checkERROR(err error) {
  if err != nil {
    revel.ERROR.Println(err)
  }
}

func checkPANIC(err error) {
  if err != nil {
    revel.ERROR.Println(err)
    panic(err)
  }
}