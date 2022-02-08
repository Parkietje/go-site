{{define "deploy"}}
    <div data-controller="deploy">
        <div class="sidebar">
            <a data-action="deploy#show" data-deploy-target-param="upload">Upload module</a>
            <a data-action="deploy#show" data-deploy-target-param="overview">List deployments</a>
        </div>
{{end}}

{{define "deploy-data"}}
        <div data-deploy-target="upload"; class="center-data hide">
            {{template "file-upload" .}}
        </div>

        <div data-deploy-target="overview"; class="center-data hide">
            <h2>Overview</h2>
            <form action="/admin/delete" method="post">
                Username:<input type="text" name="hash">
                <input type="submit" value="Delete user">
            </form>
        </div>
    </div>
{{end}}

{{define "deploy-controller"}}
    Stimulus.register("deploy", class extends Controller {

        static targets = ["upload", "overview"]

        initialize(){
          console.log('init')
        }

        // Called via:
        // data-action="click->admin#show data-admin-target-param=<your_target>
        show({ params }) {
            // show target
            let target = params["target"]
            
            if (target == "upload") {
                this.uploadTarget.style.display = "block"
                this.overviewTarget.style.display = "none"
            } else if (target == "overview"){
                this.overviewTarget.style.display = "block"
                this.uploadTarget.style.display = "none"
            }
        }     
    })
    {{template "file-upload-controller" .}}
{{end}}