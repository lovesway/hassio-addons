{{define "content"}}
<form id="addActionForm">
    <div class="form-group">
        <label for="inputDevices">Devices</label>
        <select class="custom-select" class="form-control" id="inputDevices" name="DeviceIDs" multiple="">
{{ range .Devices }}
            <option value="{{.ID}}">{{.Name}}</option>
{{ end }}
        </select>
    </div>
    <div class="form-group">
        <label for="inputCommand">Command</label>
        <select class="custom-select" class="form-control" id="inputCommand" name="Command">
            <option value=""></option>
        </select>
    </div>
    <div class="form-group">
        <label for="inputParameter">Parameter</label>
        <input type="text" class="form-control" id="inputParameter" name="Parameter" value="">
    </div>
    <div class="form-group">
        <label for="inputGlobalParameter">Use Global Parameter</label>
        <select class="custom-select" class="form-control" id="inputGlobalParameter" aria-describedby="inputGlobalParameterHelp" name="GlobalParameter">
            <option value="false">false</option>
            <option value="GlobalSpeed">GlobalSpeed</option>
            <option value="GlobalParameter1">GlobalParameter1</option>
            <option value="GlobalParameter2">GlobalParameter2</option>
        </select>
        <small id="inputGlobalParameterHelp" class="form-text text-muted">Use a Global Parameter set in the Show configuration (only works when running in show context).</small>
    </div>
    <button type="submit" class="btn btn-primary">Submit</button>
</form>
<script>
$(document).ready(function() {
    populateCommands();

    $('#addActionForm').submit(function() {
        event.preventDefault();

        if ($.trim($("#inputDevices").val()) === "" ) {
            alert('Please select at least one Device.');
            return false;
        } else
        if ($.trim($("#inputCommand").val()) === "" ) {
            alert('Please select a Command.');
            return false;
        }

        var formData = JSON.stringify($(this).serializeFormJSON());
        var objData = jQuery.parseJSON(formData);

        if (typeof objData.DeviceIDs == "string") {
            objData['DeviceIDs'] = [objData.DeviceIDs]
        }
        formData = JSON.stringify(objData);

        $.post("api/v1/scene/{{.SceneID}}/group/{{.GroupID}}/action", formData, function(data) {
            if (data.Error != false) {
                alert("Error: " + data.Message);
            } else {
                $('#focusSelector').text('#groupRow{{.GroupID}}')
                populateContent();
                $("#addActionModal").dialog("close");
            }
        }, "json");

        return false;
    });

    $("#inputDevices").change(function() {
        populateCommands($("#inputCommand").val());
    });

    function populateCommands(selected) {
        var deviceCommands = {};{{ range .DeviceTypes }}
        deviceCommands[{{.ID}}] = { {{ range .Commands }}
            "{{.Name}}": "{{.Description}}",{{ end }}
        };
{{ end }}
        var deviceToDeviceType = {};{{ range .Devices }}
        deviceToDeviceType[{{.ID}}] = {{.Type.ID}};{{ end }}
        $("#inputCommand").empty();
        // Device selection changed, need to populate the commands
        if($("#inputDevices").val()){            
            $("#inputCommand").attr('disabled',false);
            var str = $("#inputDevices").val().toString();
            var deviceId = str.split(",")[0];
            if (deviceId == "") {
                $("#inputCommand").attr('disabled',true);
                $("#inputCommand").empty();
                return;
            }
            typeId = deviceToDeviceType[deviceId];
            Object.keys(deviceCommands[typeId]).forEach(key => {
                if (key == selected) {
                    $("#inputCommand").append('<option value="'+key+'" selected>'+deviceCommands[typeId][key]+'</option>'); 
                } else {
                    $("#inputCommand").append('<option value="'+key+'">'+deviceCommands[typeId][key]+'</option>'); 
                }
            });
        } else {
            $("#inputCommand").attr('disabled',true);
            $("#inputCommand").empty();
        }
    }
});
</script>
{{end}}