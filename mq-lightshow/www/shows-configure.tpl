{{define "content"}}
<form id="showConfigureForm">
    <div class="form-group">
        <label for="inputName">Show Name</label>
        <input type="text" class="form-control" id="inputName" name="Name" value="{{.Show.Name}}">
    </div>
    <div class="form-group">
        <label for="inputTopic">MQTT Topic</label>
        <input type="text" class="form-control" id="inputTopic" aria-describedby="inputTopicHelp" name="Topic" value="{{.Show.Topic}}">
        <small id="inputTopichelp" class="form-text text-muted">If configured, show will be bound to this topic. Note that you may need to go to the MQTT page and disconnect/reconnect so that the topic will be subscribed to.</small>
    </div>
    <div class="form-group">
        <label for="inputType">Repeat</label>
        <select class="custom-select" class="form-control" id="inputRepeat" name="Repeat">
            <option value="false"{{if (eq .Show.Repeat false)}} selected{{end}}>false</option>
            <option value="true"{{if (eq .Show.Repeat true)}} selected{{end}}>true</option>
        </select>
    </div>
    <div class="form-group">
        <label for="inputGlobalDelay">Global Delay</label>
        <input type="text" class="form-control" id="inputGlobalDelay" aria-describedby="inputGlobalDelayHelp" name="GlobalDelay" value="{{.GlobalDelay}}">
        <small id="inputGlobalDelayHelp" class="form-text text-muted">This float value is exported to the scene groups in this show.</small>
    </div>
    <div class="form-group">
        <label for="inputGlobalSpeed">Global Speed</label>
        <input type="text" class="form-control" id="inputGlobalSpeed" aria-describedby="inputGlobalSpeedHelp" name="GlobalSpeed" value="{{.GlobalSpeed}}">
        <small id="inputGlobalSpeedHelp" class="form-text text-muted">This integer value is exported to the scene action parameters in this show.</small>
    </div>
    <div class="form-group">
        <label for="inputGlobalParameterValue1">Global Parameter Value 1</label>
        <input type="text" class="form-control" id="inputGlobalParameterValue1" aria-describedby="inputGlobalParameterValue1Help" name="GlobalParameter1" value="{{.Show.GlobalParameter1}}">
        <small id="inputGlobalParameterValue1Help" class="form-text text-muted">This string value is exported to the scene action parameters in this show.</small>
    </div>
    <div class="form-group">
        <label for="inputGlobalParameterValue1">Global Parameter Value 2</label>
        <input type="text" class="form-control" id="inputGlobalParameterValue2" aria-describedby="inputGlobalParameterValue2Help" name="GlobalParameter2" value="{{.Show.GlobalParameter2}}">
        <small id="inputGlobalParameterValue2Help" class="form-text text-muted">This string value is exported to the scene action parameters in this show.</small>
    </div>
    <button type="submit" class="btn btn-primary">Submit</button>
</form>
<script>
$(document).ready(function() {
    $('#showConfigureForm').submit(function(event) {
        event.preventDefault();

        if ($.trim($("#inputName").val()) === "" ) {
            alert('Please fill out the Show Name.');
            return false;
        }

        var formData = JSON.stringify($(this).serializeFormJSON());

        $.post("api/v1/show/{{.Show.ID}}/configure", formData, function(data) {
            if (data.Error != false) {
                alert ("Error: " + data.Message)
            } else {
                populateContent()
                $("#configureShowModal").dialog("close");
            }
        }, "json");

        return false;
    });
});
</script>
{{end}}