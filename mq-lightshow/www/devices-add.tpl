{{define "content"}}
<h1>Adding Device</h1>
<form method="POST" id="deviceForm">
    <div class="form-group">
        <label for="inputName">Device Name</label>
        <input type="text" class="form-control" id="inputName" aria-describedby="inputNameHelp" name="name" value="">
        <small id="inputNamehelp" class="form-text text-muted">example: Living Room Light 1</small>
    </div>
    <div class="form-group">
        <label for="inputName">MQTT Topic</label>
        <input type="text" class="form-control" id="inputTopic" name="topic" value="">
    </div>
    <div class="form-group">
        <label for="inputType">Device Type</label>
        <select class="custom-select" class="form-control" id="inputType" name="type">
            <option value=""></option>
{{ range .DeviceTypes }}
            <option value="{{.ID}}">{{.Name}}</option>
{{ end }}
        </select>
    </div>
    <button type="submit" class="btn btn-primary">Submit</button>
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