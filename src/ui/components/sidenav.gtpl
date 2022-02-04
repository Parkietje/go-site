{{define "sidenav"}}

    <div class="sidenav">
        <a onclick="show('add')">Add user</a>
        <a onclick="show('delete')">Delete user</a>
    </div>

    <div class="add">
        <h2>Add user</h2>
        <form action="/admin/add" method="post">
            Username:<input type="text" name="username">
            Password:<input type="password" name="password">
            Salt:<input type="text" name="salt">
            <input type="submit" value="Add user">
        </form>
    </div>

    <div class="delete">
        <h2>Delete user</h2>
        <form action="/admin/delete" method="post">
            Hash of username:<input type="text" name="hash">
            <input type="submit" value="Delete user">
        </form>
    </div>


    <script>
        function hideAll() {
    	    const names = ['add', 'delete']
    	    function get(thing) { return document.getElementsByClassName(thing)[0]}
    	        for (const n of names) {
    		        get(n).style.display = "none";
    	        }
        }
    </script>

    <script>
        hideAll();
        function show(name) {
  	        const names = ['add', 'delete']
  	        function get(thing) { return document.getElementsByClassName(thing)[0]}
  	        var x = get(name)
  	        if (x.style.display === "none") {
  		        x.style.display = "block";
  		        for (const n of names) {
  			        if(n == name){continue} get(n).style.display ="none" ;
  		        }
  	        } 
        }
    </script>
    
    <style>
        .add {
            display:block
        }
        .delete {
            display:block
        }

        body {
          font-family: "Lato", sans-serif;
        }   
        .sidenav {
          height: 100%;
          width: 200px;
          position: absolute;
          z-index: 1;
          top: 175px;
          left: 0;
          background-color: #34495E;
          overflow-x: hidden;
          padding-top: 20px;
        }   
        .sidenav a {
          padding: 6px 8px 6px 16px;
          text-decoration: none;
          font-size: 25px;
          color: #F7F9FA;
          display: block;
        }   
        .sidenav a:hover {
          color: #f1f1f1;
        }   
        .main {
          margin-left: 160px; /* Same as the width of the sidenav */
          font-size: 28px; /* Increased text to enable scrolling */
          padding: 0px 10px;
        }   
        @media screen and (max-height: 450px) {
          .sidenav {padding-top: 15px;}
          .sidenav a {font-size: 18px;}
        }

        nav2 ul{height:170px; width:100%;margin-bottom: 10px}
        nav2 ul{overflow:hidden; overflow-y:scroll;} 
    </style>
{{end}}