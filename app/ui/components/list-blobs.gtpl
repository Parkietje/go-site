{{define "list-blobs"}}
    <div data-controller="list-blobs">
        <label> Choose the File to upload: </label>
        <input type="text" id="text" /> <br /><br />
        <button data-action="click->list-blobs#listBlobs">Submit</button>
    </div>
{{end}}

{{define "list-blobs-controller"}}
    Stimulus.register("list-blobs", class extends Controller {
        
        listBlobs() {
            console.log("test")

            let text = document.getElementById("text").value;
            console.log(text)
            let formData = new FormData();
            formData.append('text', text)

            fetch('http://localhost:4000/deploy/list', {method: 'POST',body: formData}).then(
                (response) => {console.log(response)}
            )
        }
    })
{{end}}