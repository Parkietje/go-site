{{define "admin-menu"}}
    <div data-controller="admin-menu-controller">
        <div class="sidebar">
            <a data-action="admin-menu-controller#showAdd">Add user</a>
            <a data-action="admin-menu-controller#showDelete">Delete user</a>
        </div>
{{end}}

{{define "admin-menu-data"}}
        <div data-admin-menu-controller-target="add"; class="center-data">
            <h2>Add user</h2>
            <form action="/admin/add" method="post">
                Username:<input type="text" name="username">
                Password:<input type="password" name="password">
                <input type="submit" value="Add user">
            </form>
        </div>

        <div data-admin-menu-controller-target="delete"; class="center-data">
            <h2>Delete user</h2>
            <form action="/admin/delete" method="post">
                Username:<input type="text" name="hash">
                <input type="submit" value="Delete user">
            </form>
        </div>
    </div>
{{end}}

{{define "admin-menu-controller"}}
    Stimulus.register("admin-menu-controller", class extends Controller {
        
      
        static targets = ["add", "delete"];
  
        // Called when Stimulus create an instance of the
        // Controller class.
        initialize(){
            this._hideAll()
        }

        _hideAll() {
    	    this.addTarget.style.display = "none"
            this.deleteTarget.style.display = "none"
        }
  
        // Called when the class is connected to the HTML element.
        connect(){
            console.log('hoi')
        }

        // Called via:
        // data-action="click->admin-menu-controller#showAdd"
        showAdd() {
            this.addTarget.style.display = "block"
            this.deleteTarget.style.display = "none"
        }

        // Called via:
        // data-action="click->admin-menu-controller#showDelete"
        showDelete() {
            this.deleteTarget.style.display = "block"
            this.addTarget.style.display = "none"
        }
              
    })
{{end}}