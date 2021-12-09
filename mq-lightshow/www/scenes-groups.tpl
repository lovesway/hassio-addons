{{define "content"}}
<h1>Scene - {{.Scene.Name}}</h1>
<p>Scenes contain Action Groups which are groups of Actions executed at once. Each group has it's own time delay that determines when the next Action Group is executed.</p>
<div class="container text-right" style="padding-bottom:15px;">
  <table class="table table-borderless">
    <tr>
      <td>Scene Runtime: {{.TotalSeconds}} seconds total ({{.TotalMinutes}} minutes)</td>
      <td class="text-right">
        <button onclick="runScene({{.Scene.ID}})" class="btn btn-sm btn-primary" title="Run Scene"><div class="icon-button-execute">&nbsp;</div></button>
        <button onclick="configureModal({{.Scene.ID}})" class="btn btn-sm btn-primary" title="Configure Scene"><div class="icon-button-gear">&nbsp;</div></button>
        <button onclick="addGroupModal({{.Scene.ID}})" class="btn btn-sm btn-primary" title="Add Group to Scene">Add Group</button>
      </td>
    </tr>
  </table>
</div>
<table class="table table-borderless">
  <tbody id="groupsContainer">
  </tbody>
</table>
</div>
<div title="Configure Scene" id="configureSceneModal"></div>
<div title="Add Group" id="addGroupModal"></div>
<div title="Configure Group" id="configureGroupModal"></div>
<div title="Add Action" id="addActionModal"></div>
<div title="Edit Action" id="editActionModal"></div>
<script>
function populateContent() {
  var groupsContainer = $('#groupsContainer');
  var groups;
  var html = "";

  var getGroups = $.getJSON('api/v1/scene/{{.Scene.ID}}/groups', function (data) {
    groups = data.Data;
  });

  getGroups.done(function() {
    groupsContainer.empty();
    for (i=0; i<groups.length; i++) {
        html += `
  <tr id="groupRow${groups[i].ID}">
    <td>
    <table class="table table-borderless table-dark mb-0">
      <tr>
        <td width="20%">order: ${groups[i].Order}</td>
        <td width="20%">delay: ${groups[i].Delay}</td>
        <td width="20%">globalDelay: ${groups[i].GlobalDelay}</td>
        <td width="40%" class="text-right">
          <button onclick="runGroup({{.Scene.ID}}, ${groups[i].ID})" class="btn btn-sm btn-primary" title="Run Actions in Group">
            <div class="icon-button-execute">&nbsp;</div>
          </button>
          <button onclick="configureGroupModal({{.Scene.ID}}, ${groups[i].ID})" class="btn btn-sm btn-primary" title="Configure Group">
            <div class="icon-button-gear">&nbsp;</div>
          </button>
          <button onclick="duplicateGroup({{.Scene.ID}}, ${groups[i].ID})" class="btn btn-sm btn-warning" title="Duplicate Group">
            <div class="icon-button-duplicate">&nbsp;</div>
          </button>
          <button onclick="deleteGroup({{.Scene.ID}}, ${groups[i].ID})" class="btn btn-sm btn-danger" title="Delete Group">
            <div class="icon-button-delete">&nbsp;</div>
          </button>
        </td>
      </tr>
    </table>
    <div id="containerActionsForGroup${groups[i].ID}" class="container hideWhenSortingGroups">
      <table class="table hideWhenSortingGroups">
        <thead>
          <tr class="thead-light">
            <th scope="col" style="display:none">#</th>
            <th scope="col">Cmd</th>
            <th scope="col">Param</th>
            <th scope="col">GParam</th>
            <th scope="col">Devices</th>
            <th scope="col"><button onclick="addActionModal({{.Scene.ID}}, ${groups[i].ID})" class="btn btn-sm btn-primary" title="Add Action to Group">Add Action</button></th>
          </tr>
        </thead>
        <tbody id="groupActions${groups[i].ID}">
        </tbody>
      </table>
    </div>
    </td>
  </tr>
  `;
    }
    groupsContainer.append(html);

    var sceneID = {{.Scene.ID}};

    var group_start_pos;
    var group_index;
    $('#groupsContainer').sortable({
        start: function(event, ui) {
          group_start_pos = ui.item.index();
          ui.item.data('start_pos', group_start_pos);
          $(".hideWhenSortingGroups").hide();
          $('#focusSelector').text('#groupsContainer');
          focusPage();
        },
        change: function(event, ui) {
            group_start_pos = ui.item.data('start_pos');
            group_index = ui.placeholder.index();
        },
        update: function(event, ui) {
          groupID = ui.item[0].id.substring(8);
          if (group_start_pos < group_index) {
            sort = group_index;
          } else {
            sort = group_index+1;
          }

          $(".hideWhenSortingGroups").show();
            $.post("api/v1/scene/"+sceneID+"/group/"+groupID+"/sort/"+sort, function(data) {
                if (data.Error != false) {
                    alert("Error: " + data.Message);
                } else {
                    $('#focusSelector').text('#groupRow'+groupID);
                    populateContent();
                }
            }, "json");
        }
      });

    for (i=0; i<groups.length; i++) {
      var actions;
      var getActions = $.getJSON("api/v1/scene/{{.Scene.ID}}/group/" + groups[i].ID + "/actions", function (data) {
        actions = data.Data;
      });
      getActions.done(function() {
        for (i=0; i<actions.length; i++) {
            thisHtml = `
          <tr id="actionRow${actions[i].ID}">
            <td style="display:none">${actions[i].ID}</td>
            <td class="">${actions[i].Command}</td> 
            <td>${actions[i].Parameter}</td>
            <td>${actions[i].GlobalParameter}</td>
            <td>`;
            looped = false
            for (ii=0; ii<actions[i].Devices.length; ii++) {
              if (looped == true) {
                 thisHtml += ", ";
              } else {
                looped = true;
              }

              thisHtml += actions[i].Devices[ii].Name;
            }
            thisHtml +=`</td>
            <td>
              <button onclick="runAction({{.Scene.ID}}, ${actions[i].GroupID}, ${actions[i].ID})" class="btn btn-sm btn-primary" title="Run Action"><div class="icon-button-execute">&nbsp;</div></button>
              <button onclick="editActionModal({{.Scene.ID}}, ${actions[i].GroupID}, ${actions[i].ID})" class="btn btn-sm btn-primary" title="Edit Action"><div class="icon-button-edit">&nbsp;</div></button>
              <button onclick="duplicateAction({{.Scene.ID}}, ${actions[i].GroupID}, ${actions[i].ID})" class="btn btn-sm btn-warning" title="Duplicate Action"><div class="icon-button-duplicate">&nbsp;</div></button>
              <button onclick="deleteAction({{.Scene.ID}}, ${actions[i].GroupID}, ${actions[i].ID})" class="btn btn-sm btn-danger" title="Delete Action"><div class="icon-button-delete">&nbsp;</div></button>
            </td>
          </tr>`;
          $('#groupActions'+ actions[i].GroupID).append(thisHtml);
        }

        focusPage();
      });
    }


    for (i=0; i<groups.length; i++) {
      var start_pos;
      var index; 

      $('#groupActions' + groups[i].ID).sortable({
        start: function(event, ui) {
            start_pos = ui.item.index();
            ui.item.data('start_pos', start_pos);
            $('#' + event.target.id +' tr').addClass('bg-light');
        },
        change: function(event, ui) {
            start_pos = ui.item.data('start_pos');
            index = ui.placeholder.index();
        },
        update: function(event, ui) {
            $('#' + event.target.id +' tr').removeClass('bg-light');
            groupID = event.target.id.substring(12);

            if (start_pos < index) {
              actionID = $('#'+event.target.id+' tr:eq('+(index-1)+') td:first-child').html();
              sort = index;
            } else {
              actionID = $('#'+event.target.id+' tr:eq('+index+') td:first-child').html();
              sort = index+1;
            }

            $.post("api/v1/scene/"+sceneID+"/group/"+groupID+"/action/"+actionID+"/sort/"+sort, function(data) {
                if (data.Error != false) {
                    alert("Error: " + data.Message);
                } else {
                    $('#focusSelector').text('#actionRow'+actionID);
                    populateContent();
                }
            }, "json");
        }
      });
    }

    focusPage();
  });

  groupsContainer.text('Loading Scene Groups from the API...');
}

function configureModal(sceneID) {
  $('#focusSelector').text('');
  $("#configureSceneModal").dialog("open");
  $.get("scenes-configure?sceneID=" + sceneID, function(html){ 
    $('#configureSceneModal').append(html);
    $('#inputName').trigger('focus');
  });
}

function runScene(sceneID) {
  $('#focusSelector').text('');
  $.post("api/v1/scene/"+sceneID+"/run", function(data) {
      if (data.Error != false) {
          alert("Error: " + data.Message);
      } else {
          $('#focusSelector').text('');
          populateContent();
      }
  }, "json");
}

function addGroupModal(sceneID) {
  $('#focusSelector').text('');
  $("#addGroupModal").dialog("open");
  $.get("scenes-groups-add?sceneID="+sceneID, function(html){ 
    $('#addGroupModal').append(html);
    $('#inputDelay').trigger('focus');
  });

}

function configureGroupModal(sceneID, groupID) {
  $('#focusSelector').text('');
  $("#configureGroupModal").dialog("open");
  $.get("scenes-groups-configure?sceneID=" + sceneID + "&groupID=" + groupID, function(html){ 
    $('#configureGroupModal').append(html);
    $('#inputDelay').trigger('focus');
  });
}

function runGroup(sceneID, groupID) {
  $('#focusSelector').text('');
  $.post("api/v1/scene/"+sceneID+"/group/"+groupID+"/run", function(data) {
      if (data.Error != false) {
          alert("Error: " + data.Message);
      } else {
          $('#focusSelector').text('#groupRow'+groupID);
          populateContent();
      }
  }, "json");
}

function duplicateGroup(sceneID, groupID) {
  $('#focusSelector').text('');
  if (confirm('Are you sure you want to duplicate this group?')) {
    $.post("api/v1/scene/"+sceneID+"/group/"+groupID+"/duplicate", function(data) {
        if (data.Error != false) {
            alert("Error: " + data.Message);
        } else {
          $('#focusSelector').text('#groupRow'+groupID);
            populateContent();
        }
    }, "json");
  }
}

function deleteGroup(sceneID, groupID) {
  $('#focusSelector').text('');
  if (confirm('Are you sure you want to delete this group along with its actions?')) {
    $.post("api/v1/scene/"+sceneID+"/group/"+groupID+"/delete", function(data) {
        if (data.Error != false) {
            alert("Error: " + data.Message);
        } else {
            $('#focusSelector').text('');
            populateContent();
        }
    }, "json");
  }
}

function addActionModal(sceneID, groupID) {
  $('#focusSelector').text('');
  $("#addActionModal").dialog("open");
  $.get("scenes-groups-actions-add?sceneID="+sceneID+"&groupID="+groupID, function(html){ 
    $('#addActionModal').append(html);
    $('#inputDevices').trigger('focus');
  });
}

function editActionModal(sceneID, groupID, actionID) {
  $('#focusSelector').text('');
  $("#editActionModal").dialog("open");
  $.get("scenes-groups-actions-edit?sceneID="+sceneID+"&groupID="+groupID+"&actionID="+actionID, function(html){ 
    $('#editActionModal').append(html);
    $('#inputDevices').trigger('focus');
  });
}

function runAction(sceneID, groupID, actionID) {
  $('#focusSelector').text('');
  $.post("api/v1/scene/"+sceneID+"/group/"+groupID+"/action/"+actionID+"/run", function(data) {
      if (data.Error != false) {
          alert("Error: " + data.Message);
      } else {
          $('#focusSelector').text('#actionRow'+actionID);
          populateContent();
      }
  }, "json");
}

function deleteAction(sceneID, groupID, actionID) {
  $('#focusSelector').text('');
  if (confirm('Are you sure you want to delete this action?')) {
    $.post("api/v1/scene/"+sceneID+"/group/"+groupID+"/action/"+actionID+"/delete", function(data) {
        if (data.Error != false) {
            alert("Error: " + data.Message);
        } else {
            $('#focusSelector').text('#groupRow'+groupID)
            populateContent();
        }
    }, "json");
  }
}

function duplicateAction(sceneID, groupID, actionID) {
  $('#focusSelector').text('');
  if (confirm('Are you sure you want to duplicate this action?')) {
    $.post("api/v1/scene/"+sceneID+"/group/"+groupID+"/action/"+actionID+"/duplicate", function(data) {
        if (data.Error != false) {
            alert("Error: " + data.Message);
        } else {
            $('#focusSelector').text('#groupRow'+groupID)
            populateContent();
        }
    }, "json");
  }
}

$(document).ready(function() {
  populateContent();
  makeDialog("#configureSceneModal");
  makeDialog("#addGroupModal")
  makeDialog("#configureGroupModal")
  makeDialog("#addActionModal")
  makeDialog("#editActionModal")
});
</script>
{{end}}
