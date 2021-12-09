{{define "content"}}
<h1>Editing Device</h1>
<form method="POST" id="deviceForm">
    <div class="form-group">
        <label for="inputName">Device Name</label>
        <input type="text" class="form-control" id="inputName" name="name" value="{{ .Device.Name }}">
    </div>
    <div class="form-group">
        <label for="inputName">MQTT Topic</label>
        <input type="text" class="form-control" id="inputTopic" name="topic" value="{{ .Device.Topic }}">
    </div>
    <div class="form-group">
        <label for="inputType">Device Type</label>
        <select class="custom-select" class="form-control" id="inputType" name="type">
            <option value=""></option>
{{ range $key, $device := .DeviceTypes }}
        <option value="{{.ID}}"{{if eq $device.Name .Name }} selected{{end}}>{{.Name}}</option>
{{ end }}
        </select>
    </div>
    <button type="submit" class="btn btn-primary">Update Device</button>
</form>
<script>
$(document).ready(function() {
    $('#deviceForm').submit(function() {
        if ($.trim($("#inputName").val()) === "" ) {
            alert('Please fill out the Device Name.');
            return false;
        } else
        if ($.trim($("#inputTopic").val()) === "" ) {
            alert('Please fill out the Device Topic.');
            return false;
        } else
        if ($.trim($("#inputType").val()) === "") {
            alert('Please select a Device Type.');
            return false;
        }
    });
});
</script>
{{end}}