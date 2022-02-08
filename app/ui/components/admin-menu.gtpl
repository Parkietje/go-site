{{define "admin"}}
    <div data-controller="admin">
        <div class="sidebar">
            <a data-action="admin#show" data-admin-target-param="add">Add user</a>
            <a data-action="admin#show" data-admin-target-param="delete">Delete user</a>
        </div>
{{end}}

{{define "admin-data"}}
        <div data-admin-target="add"; class="center-data hide">
            <h2>Add user</h2>
            <form action="/admin/add" method="post">
                Username:<input type="text" name="username">
                Password:<input type="password" name="password">
                <input type="submit" value="Add user">
            </form>
        </div>

        <div data-admin-target="delete"; class="center-data hide">
            <h2>Delete user</h2>
            <form action="/admin/delete" method="post">
                Username:<input type="text" name="hash">
                <input type="submit" value="Delete user">
            </form>
        </div>
    </div>
{{end}}

{{define "admin-controller"}}
    Stimulus.register("admin", class extends Controller {

        static targets = ["add", "delete"]

        initialize(){
          console.log('init')
        }

        // Called via:
        // data-action="click->admin#show data-admin-target-param=<your_target>
        show({ params }) {
            // show target
            let target = params["target"]
            
            if (target == "add") {
                this.addTarget.style.display = "block"
                this.deleteTarget.style.display = "none"
            } else if (target == "delete"){
                this.deleteTarget.style.display = "block"
                this.addTarget.style.display = "none"
            }
        }     
    })
{{end}}