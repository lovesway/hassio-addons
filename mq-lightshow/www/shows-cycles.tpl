{{define "content"}}
<h1>Light Show - {{.Show.Name}}</h1>
<p>Each Scene Cycle contains a Scene along with settings for the Scene's cycle.</p>
<table class="table">
  <thead class="thead-dark">
    <tr>
      <th scope="col">Scene</th>
      <th scope="col">Cycles</th>
      <th scope="col">End Delay</th>
      <th scope="col">Loop Include</th>
      <th scope="col" title="Global Delay">GD</th>
      <th scope="col" title="Global Speed">GS</th>
      <th scope="col" title="Global Parameter 1">GP1</th>
      <th scope="col" title="Global Parameter 2">GP2</th>
      <th scope="col" class="text-right">
      {{if eq .Show.Running true}}
        <button onclick="stopShow({{.Show.ID}})" class="btn btn-sm btn-danger" title="Stop Show">
          <div class="icon-button-execute">&nbsp;</div>
        </button>
      {{else}}
        <button onclick="startShow({{.Show.ID}})" class="btn btn-sm btn-primary" title="Start Show">
          <div class="icon-button-execute">&nbsp;</div>
        </button>
      {{end}}
        <button onclick="configureModal({{.Show.ID}})" class="btn btn-sm btn-primary" title="Configure Show"><div class="icon-button-gear">&nbsp;</div></button>
        <button onclick="addModal({{.Show.ID}})" class="btn btn-sm btn-primary" title="Add a New Scene Cycle">Add Cycle</button>
      </th>
    </tr>
  </thead>
  <tbody id="cyclesContainer">
  </tbody>
</table>
<div title="Add Cycle" id="addCycleModal"></div>
<div title="Edit Cycle" id="editCycleModal"></div>
<div title="Configure Show" id="configureShowModal"></div>
<script>
function populateContent() {
  var cyclesContainer = $('#cyclesContainer');

  $.getJSON('api/v1/show/{{.Show.ID}}/cycles', function (data) {
    cyclesContainer.empty();
    cycles = data.Data;

    var html = "";

    for (i=0; i<cycles.length; i++) {
      if (cycles[i].GlobalDelay == 0) {
        var gDelay = '';
      } else {
        var gDelay = cycles[i].GlobalDelay;
      }

      if (cycles[i].GlobalSpeed == 0) {
        var gSpeed = '';
      } else {
        var gSpeed = cycles[i].GlobalSpeed;
      }

      html +=`
      <tr>
        <td>${cycles[i].Scene.Name}</td>
        <td>${cycles[i].SceneCycles}</td>
        <td>${cycles[i].EndDelay}</td>
        <td>${cycles[i].LoopInclude}</td>
        <td>${gDelay}</td>
        <td>${gSpeed}</td>
        <td>${cycles[i].GlobalParameter1}</td>
        <td>${cycles[i].GlobalParameter2}</td>
        <td>
          <button onclick="editModal(${cycles[i].ID})" class="btn btn-sm btn-primary" title="Edit Scene Cycle">Edit Cycle</button>
          <button onclick="deleteShowCycle({{.Show.ID}}, ${cycles[i].ID})" class="btn btn-sm btn-danger" title="Delete Cycle"><div class="icon-button-delete">&nbsp;</div></button>
        </td>
      </tr>`;
      }

    cyclesContainer.append(html);
  });

  cyclesContainer.html('<tr><td colspan="9">Loading Cycles from the API...</td></tr>');
}

function configureModal(id) {
  $("#configureShowModal").dialog("open");
  $.get("shows-configure?showID="+id, function(html){
    $('#configureShowModal').append(html);
    $('#inputName').trigger('focus');
  });

}

function addModal(showID) {
  $("#addCycleModal").dialog("open");
  $.get("shows-cycles-add?showID=" + showID, function(html){ 
    $('#addCycleModal').append(html);
    $('#inputScene').trigger('focus');
  });
}

function editModal(cycleID) {
  $("#editCycleModal").dialog("open");
  $.get("shows-cycles-edit?cycleID=" + cycleID, function(html){ 
    $('#editCycleModal').append(html);
    $('#inputScene').trigger('focus');
  });
}

function startShow(id) {
  $.post("api/v1/show/"+id+"/start", function(data) {
      if (data.Error != false) {
          alert ("Error: " + data.Message)
      } else {
          location.reload(true);
      }
  }, "json");
}

function stopShow(id) {
  $.post("api/v1/show/"+id+"/stop", function(data) {
      if (data.Error != false) {
          alert ("Error: " + data.Message)
      } else {
          location.reload(true);
      }
  }, "json");
}

function deleteShowCycle(showID, cycleID) {
  if (confirm('Are you sure you want to delete this cycle?')) {
    $.post("api/v1/show/" + showID + "/cycle/" + cycleID + "/delete", function(data) {
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
  makeDialog("#addCycleModal");
  makeDialog("#editCycleModal");
  makeDialog("#configureShowModal");
});
</script>
{{end}}
