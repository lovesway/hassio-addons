{{define "content"}}
<h1>Light Shows</h1>
<p>Shows are collections of Scenes with individual settings. You must create at least one scene before you can create a show.</p>
<table class="table">
  <thead class="thead-dark">
    <tr>
      <th scope="col">Name</th>
      <th scope="col">Repeat</th>
      <th scope="col">GD</th>
      <th scope="col">GS</th>
      <th scope="col">GP1</th>
      <th scope="col">GP2</th>
      <th scope="col" class="text-right"><button class="btn btn-sm btn-primary" onclick="addModal()" title="Add a New Light Show">Add Show</button></th>
    </tr>
  </thead>
  <tbody id="showsContainer">
  </tbody>
</table>
<div title="Add Show" id="addShowModal"></div>
<div title="Configure Show" id="configureShowModal"></div>
<script>
function populateContent() {
   var showsContainer = $('#showsContainer');

    $.getJSON('api/v1/shows', function (data) {
      showsContainer.empty();
      shows = data.Data;
      var html = "";

      for (i=0; i<shows.length; i++) {
        if (shows[i].GlobalDelay == 0) {
          var gDelay = '';
        } else {
          var gDelay = shows[i].GlobalDelay;
        }

        if (shows[i].GlobalSpeed == 0) {
          var gSpeed = '';
        } else {
          var gSpeed = shows[i].GlobalSpeed;
        }

        html +=`
    <tr>
      <td>${shows[i].Name}</td>
      <td>${shows[i].Repeat}</td>
      <td>${gDelay}</td>
      <td>${gSpeed}</td>
      <td>${shows[i].GlobalParameter1}</td>
      <td>${shows[i].GlobalParameter2}</td>
      <td>`;

        if (shows[i].Running == true) {
        html +=`
        <button onclick="stopShow(${shows[i].ID})" class="btn btn-sm btn-danger" title="Stop Show">
          <div class="icon-button-execute">&nbsp;</div>
        </button>`;
        } else {
        html +=`
        <button onclick="startShow(${shows[i].ID})" class="btn btn-sm btn-primary" title="Start Show">
          <div class="icon-button-execute">&nbsp;</div>
        </button>`;
        }

        html +=`
        <button onclick="configureModal(${shows[i].ID})" class="btn btn-sm btn-primary" title="Configure Show"><div class="icon-button-gear">&nbsp;</div></button>
        <a href="shows-cycles?showID=${shows[i].ID}" class="btn btn-sm btn-primary" title="Edit Show"><div class="icon-button-edit">&nbsp;</div></a>
        <button onclick="deleteShow(${shows[i].ID})" class="btn btn-sm btn-danger" title="Delete Show"><div class="icon-button-delete">&nbsp;</div></button>
      </td>
    </tr>`;
      }

        showsContainer.append(html);
    });

    showsContainer.html('<tr><td colspan="7">Loading Shows from the API...</td></tr>');
}

function addModal() {
  $("#addShowModal").dialog("open");
  $.get("shows-add", function(html){ 
    $('#addShowModal').append(html);
    $('#inputName').trigger('focus');
  });
}

function configureModal(id) {
  $("#configureShowModal").dialog("open");
  $.get("shows-configure?showID="+id, function(html){
    $('#configureShowModal').append(html);
    $('#inputName').trigger('focus');
  });
}

function startShow(id) {
  $.post("api/v1/show/"+id+"/start", function(data) {
      if (data.Error != false) {
          alert ("Error: " + data.Message)
      } else {
          populateContent()
      }
  }, "json");
}

function stopShow(id) {
  $.post("api/v1/show/"+id+"/stop", function(data) {
      if (data.Error != false) {
          alert ("Error: " + data.Message)
      } else {
          populateContent()
      }
  }, "json");
}

function deleteShow(id) {
  if (confirm('Are you sure you want to delete this show?')) {
    $.post("api/v1/show/"+id+"/delete", function(data) {
        if (data.Error != false) {
            alert ("Error: " + data.Message)
        } else {
            populateContent()
        }
    }, "json");
  }  
}

$(document).ready(function() {
  populateContent();
  makeDialog("#addShowModal");
  makeDialog("#configureShowModal");
});
</script>
{{end}}
