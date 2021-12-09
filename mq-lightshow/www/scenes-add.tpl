{{define "content"}}
<form method="POST" id="sceneAddForm">
    <div class="form-group">
        <label for="inputName">Scene Name</label>
        <input type="text" class="form-control" id="inputName" name="Name" value="">
    </div>
    <div class="form-group">
        <label for="inputAllowedDevices">Allowed Devices</label> <button class="btn btn-sm btn-info" title="Clear Selection" onclick="return ResetDevices()">clear selection</button>
        <select class="custom-select" class="form-control" id="inputAllowedDevices" aria-describedby="inputAllowedDevicesHelp" name="AllowedDeviceIDs" multiple="">
{{ range .Devices }}
            <option value="{{.ID}}">{{.Name}}</option>
{{ end }}
        </select>
        <small id="inputAllowedDevicesHelp" class="form-text text-muted">Limit the devices shown when adding actions in this scene.</small>
    </div>
    <button type="submit" value="submit" class="btn btn-primary">Submit</button>
</form>
<script>
function ResetDevices() {
    $('#inputAllowedDevices').val('');
    return false;
}

$(document).ready(function() {
    $('#sceneAddForm').submit(function(event) {
        event.preventDefault();

        if ($.trim($("#inputName").val()) === "" ) {
            alert('Please fill out the Scene Name.');
            return false;
        }

        var formData = JSON.stringify($(this).serializeFormJSON());
        var objData = jQuery.parseJSON(formData);

        if (typeof objData.AllowedDeviceIDs == "string") {
            objData['AllowedDeviceIDs'] = [objData.AllowedDeviceIDs]
        }
        formData = JSON.stringify(objData);

        $.post("api/v1/scene", formData, function(data) {
            if (data.Error != false) {
                alert("Error: " + data.Message);
            } else {
                populateContent();
                $("#addSceneModal").dialog("close");
            }
        }, "json");

        return false;
    });
});
</script>
{{end}}