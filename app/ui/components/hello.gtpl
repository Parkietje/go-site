{{define "hello"}}
    <div data-controller="hello">
        <button data-action="click->hello#greet">Greet</button>
    </div>
{{end}}

{{define "hello-controller"}}
    Stimulus.register("hello", class extends Controller {
      greet() {
        console.log("Hello, Stimulus!", this.element)
      }
    })
{{end}}