{{define "content"}}
<form method="POST" id="sceneConfigureForm">
    <div class="form-group">
        <label for="inputName">Scene Name</label>
        <input type="text" class="form-control" id="inputName" name="Name" value="{{ .Scene.Name }}">
    </div>
    <div class="form-group">
        <label for="inputType">Allowed Devices</label> <button class="btn btn-sm btn-info" title="Clear Selection" onclick="return ResetDevices()">clear selection</button>
        <select class="custom-select" class="form-control" id="inputDevices" aria-describedby="inputDevicesHelp" name="AllowedDeviceIDs" multiple="">
{{ range .Devices }}
            <option value="{{.ID}}"{{if .Selected}} selected{{end}}>{{.Name}}</option>
{{ end }}
        </select>
        <small id="inputNameHelp" class="form-text text-muted">Limit the devices shown when adding actions in this scene.</small>
    </div>
    <button type="submit" class="btn btn-primary">Update Scene</button>
</form>
<script>
function ResetDevices() {
    $('#inputDevices').val('');
    return false;
}
$(document).ready(function() {
    $('#sceneConfigureForm').submit(function(event) {
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

        $.post("api/v1/scene/{{.Scene.ID}}/configure", formData, function(data) {
            if (data.Error != false) {
                alert("Error: " + data.Message);
            } else {
                populateContent();
                $("#configureSceneModal").dialog("close");
            }
        }, "json");

        return false;
    });
});
</script>
{{end}}