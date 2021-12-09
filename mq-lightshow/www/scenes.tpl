{{define "content"}}
<h1>Scenes</h1>
<p>Scenes are reusable building blocks for routines.</p>
<table class="table">
  <thead class="thead-dark">
    <tr>
      <th scope="col">Name</th>
      <th scope="col" class="text-right"><button onclick="addModal()" class="btn btn-sm btn-primary" title="Add a New Scene">Add Scene</button></th>
    </tr>
  </thead>
  <tbody id="scenesContainer">
  </tbody>
</table>
<div title="Add Scene" id="addSceneModal"></div>
<div title="Configure Scene" id="configureSceneModal"></div>
<script>
function populateContent() {
   var scenesContainer = $('#scenesContainer');

    $.getJSON('api/v1/scenes', function (data) {
      scenesContainer.empty();
      scenes = data.Data;
      var html = "";

      for (i=0; i<scenes.length; i++) {
      html +=`
    <tr id="sceneRow${scenes[i].ID}">
      <td>${scenes[i].Name}</td>
      <td class="text-right">
        <button onclick="runScene(${scenes[i].ID})" class="btn btn-sm btn-primary" title="Run Scene"><div class="icon-button-execute">&nbsp;</div></button>
        <button onclick="configureModal(${scenes[i].ID})" class="btn btn-sm btn-primary" title="Configure Scene"><div class="icon-button-gear">&nbsp;</div></button>
        <a href="scenes-groups?sceneID=${scenes[i].ID}" class="btn btn-sm btn-primary" title="Edit Scene"><div class="icon-button-edit">&nbsp;</div></a>
        <button onclick="duplicateScene(${scenes[i].ID})" class="btn btn-sm btn-warning" title="Duplicate Scene"><div class="icon-button-duplicate">&nbsp;</div></button>
        <button onclick="deleteScene(${scenes[i].ID})" class="btn btn-sm btn-danger" title="Delete Scene"><div class="icon-button-delete">&nbsp;</div></button>
      </td>
    </tr>`;
      }
      scenesContainer.append(html);
    });

    scenesContainer.html('<tr><td colspan="2">Loading Scenes from the API...</td></tr>');
}

function addModal() {
  $("#addSceneModal").dialog("open");
  $.get("scenes-add", function(html){ 
    $('#addSceneModal').append(html);
    $('#inputName').trigger('focus');
  });
}

function configureModal(sceneID) {
  $("#configureSceneModal").dialog("open");
  $.get("scenes-configure?sceneID=" + sceneID, function(html){ 
    $('#configureSceneModal').append(html);
    $('#inputName').trigger('focus');
  });
}

function runScene(sceneID) {
  $.post("api/v1/scene/"+sceneID+"/run", function(data) {
      if (data.Error != false) {
          alert ("Error: " + data.Message)
      } else {
          populateContent()
      }
  }, "json");
}

function deleteScene(sceneID) {
  if (confirm('Are you sure you want to delete this scene?')) {
    $.post("api/v1/scene/"+sceneID+"/delete", function(data) {
        if (data.Error != false) {
            alert ("Error: " + data.Message)
        } else {
            populateContent()
        }
    }, "json");
  }
}

function duplicateScene(sceneID) {
  $('#focusSelector').text('');
  if (confirm('Are you sure you want to duplicate this scene?')) {
    $.post("api/v1/scene/"+sceneID+"/duplicate", function(data) {
        if (data.Error != false) {
            alert("Error: " + data.Message);
        } else {
          $('#focusSelector').text('#sceneRow'+sceneID);
            populateContent();
        }
    }, "json");
  }
}

$(document).ready(function() {
  populateContent();
  makeDialog("#addSceneModal");
  makeDialog("#configureSceneModal");
});
</script>
{{end}}
