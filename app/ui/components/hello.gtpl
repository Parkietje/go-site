{{define "hello"}}
    <div data-controller="hello">
        <button data-action="click->hello#greet">Greet</button>
    </div>
{{end}}

{{define "hello-controller"}}
    Stimulus.register("hello", class extends Controller {

      // Called when Stimulus create an instance of the
      // Controller class.
      initialize(){
          console.log('init')
      }

      // Called when the class is connected to the HTML element.
      connect(){
          console.log('connected')
      }

      // Called via:
      // data-action="click->hello-controller#greet"
      greet() {
        console.log("Hello, Stimulus!", this.element)
      }
    })
{{end}}