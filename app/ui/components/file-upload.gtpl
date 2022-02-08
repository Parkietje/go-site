{{define "file-upload"}}
    <div data-controller="file-upload">
        <label> Choose the File to upload: </label>
        <input type="file" id="file-to-upload" /> <br /><br />
        <button data-action="click->file-upload#getFile">Submit</button>
    </div>
{{end}}

{{define "file-upload-controller"}}
    Stimulus.register("file-upload", class extends Controller {
        
        getFile() {
            console.log("test")
            var fullPath = document.getElementById('file-to-upload').value;
            
            console.log(fullPath)

            let files = document.getElementById("file-to-upload").files;
            let formData = new FormData();

            // Append files to files array
            for (let i = 0; i < files.length; i++) {
                let file = files[i]
                console.log(file)
                formData.append('files[]', file)
            }

            fetch('http://localhost:4000/deploy', {method: 'POST',body: formData}).then(
                (response) => {console.log(response)}
            )
        }
    })
{{end}}