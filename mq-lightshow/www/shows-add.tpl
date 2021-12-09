{{define "content"}}
<form method="POST" id="showAddForm">
    <div class="form-group">
        <label for="inputName">Show Name</label>
        <input type="text" class="form-control" id="inputName" aria-describedby="inputNameHelp" name="Name" value="">
        <small id="inputNamehelp" class="form-text text-muted">example: Porch Slow Fade</small>
    </div>
    <div class="form-group">
        <label for="inputTopic">MQTT Topic</label>
        <input type="text" class="form-control" id="inputTopic" aria-describedby="inputTopicHelp" name="Topic" value="">
        <small id="inputTopichelp" class="form-text text-muted">If configured, show will be bound to this topic. Note that you may need to go to the MQTT page and disconnect/reconnect so that the topic will be subscribed to.</small>
    </div>
    <div class="form-group">
        <label for="inputRepeat">Repeat</label>
        <select class="custom-select" class="form-control" id="inputRepeat" aria-describedby="inputRepeatHelp" name="Repeat">
            <option value="false">false</option>
            <option value="true">true</option>
        </select>
        <small id="inputRepeatHelp" class="form-text text-muted">Run this show in a loop.</small>
    </div>
    <div class="form-group">
        <label for="inputGlobalDelay">Global Delay</label>
        <input type="text" class="form-control" id="inputGlobalDelay" aria-describedby="inputGlobalDelayHelp" name="GlobalDelay" value="">
        <small id="inputGlobalDelayHelp" class="form-text text-muted">This float value is exported to the scene groups in this show.</small>
    </div>
    <div class="form-group">
        <label for="inputGlobalSpeed">Global Speed</label>
        <input type="text" class="form-control" id="inputGlobalSpeed" aria-describedby="inputGlobalSpeedHelp" name="GlobalSpeed" value="">
        <small id="inputGlobalSpeedHelp" class="form-text text-muted">This integer value is exported to the scene action parameters in this show.</small>
    </div>
    <div class="form-group">
        <label for="inputGlobalParameterValue1">Global Parameter Value 1</label>
        <input type="text" class="form-control" id="inputGlobalParameterValue1" aria-describedby="inputGlobalParameterValue1Help" name="GlobalParameter1" value="">
        <small id="inputGlobalParameterValue1Help" class="form-text text-muted">This string value is exported to the scene action parameters in this show.</small>
    </div>
    <div class="form-group">
        <label for="inputGlobalParameterValue1">Global Parameter Value 2</label>
        <input type="text" class="form-control" id="inputGlobalParameterValue2" aria-describedby="inputGlobalParameterValue2Help" name="GlobalParameter2" value="">
        <small id="inputGlobalParameterValue2Help" class="form-text text-muted">This string value is exported to the scene action parameters in this show.</small>
    </div>
    <button type="submit" class="btn btn-primary">Submit</button>
</form>
<script>
$(document).ready(function() {
    $('#showAddForm').submit(function(event) {
        event.preventDefault();

        if ($.trim($("#inputName").val()) === "" ) {
            alert('Please fill out the Show Name.');
            return false;
        }

        var formData = JSON.stringify($(this).serializeFormJSON());

        $.post("api/v1/show", formData, function(data) {
            if (data.Error != false) {
                alert ("Error: " + data.Message)
            } else {
                populateContent()
                $("#addShowModal").dialog("close");
            }
        }, "json");

        return false;
    });
});
</script>
{{end}}