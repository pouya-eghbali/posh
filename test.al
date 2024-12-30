from "fmt" import * as fmt

fn main(name string) {
  message = io.Format("Hello, %s!", name)
  
  result =
    message
      | echo()
      | tr("[:lower:]", "[:upper:]")
      | lolcat(-f)

  io.Line(result)
}
