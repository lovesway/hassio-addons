{{define "content"}}
<form id="cycleForm">
    <div class="form-group">
        <label for="inputScene">Scene</label>
        <select class="custom-select" class="form-control" id="inputScene" name="SceneID">
            <option value=""></option>
{{ range .Scenes }}
            <option value="{{.ID}}">{{.Name}}</option>
{{ end }}
        </select>
    </div>
    <div class="form-group">
        <label for="inputCycles">Scene Cycles</label>
        <input type="text" class="form-control" id="inputCycles" aria-describedby="inputCyclesHelp" name="SceneCycles" value="1">
        <small id="inputCyclesHelp" class="form-text text-muted">How many times to run the scene.</small>
    </div>
    <div class="form-group">
        <label for="inputEndDelay">End Delay</label>
        <input type="text" class="form-control" id="inputEndDelay" aria-describedby="inputEndDelayHelp" name="EndDelay" value="0.0">
        <small id="inputEndDelayHelp" class="form-text text-muted">Optional time in seconds to delay after last cycle is complete.</small>
    </div>
    <div class="form-group">
        <label for="inputLoopInclude">Loop Include</label>
        <select class="custom-select" class="form-control" id="inputLoopInclude" aria-describedby="inputLoopIncludeHelp" name="LoopInclude">
            <option value="true">true</option>
            <option value="false">false</option>
        </select>
        <small id="inputLoopIncludeHelp" class="form-text text-muted">Whether or not to include this Scene when looping (if false, it will only be executed once). This is handy for setup scenes you only need to run once.</small>
    </div>
    <div class="form-group">
        <label for="inputGlobalDelay">Global Delay</label>
        <input type="text" class="form-control" id="inputGlobalDelay" aria-describedby="inputGlobalDelayHelp" name="GlobalDelay" value="">
        <small id="inputGlobalDelayHelp" class="form-text text-muted">This float value is exported to the scene groups in this show. If empty, the Global values defined in the Show are used.</small>
    </div>
    <div class="form-group">
        <label for="inputGlobalSpeed">Global Speed</label>
        <input type="text" class="form-control" id="inputGlobalSpeed" aria-describedby="inputGlobalSpeedHelp" name="GlobalSpeed" value="">
        <small id="inputGlobalSpeedHelp" class="form-text text-muted">This integer value is exported to the scene action parameters in this show. If empty, the Global values defined in the Show are used.</small>
    </div>
    <div class="form-group">
        <label for="inputGlobalParameterValue1">Global Parameter Value 1</label>
        <input type="text" class="form-control" id="inputGlobalParameterValue1" aria-describedby="inputGlobalParameterValue1Help" name="GlobalParameter1" value="">
        <small id="inputGlobalParameterValue1Help" class="form-text text-muted">This string value is exported to the scene action parameters in this show. If empty, the Global values defined in the Show are used.</small>
    </div>
    <div class="form-group">
        <label for="inputGlobalParameterValue2">Global Parameter Value 2</label>
        <input type="text" class="form-control" id="inputGlobalParameterValue2" aria-describedby="inputGlobalParameterValue2Help" name="GlobalParameter2" value="">
        <small id="inputGlobalParameterValue2Help" class="form-text text-muted">This string value is exported to the scene action parameters in this show. If empty, the Global values defined in the Show are used.</small>
    </div>
    <button type="submit" class="btn btn-primary">Submit</button>
</form>
<script>
$(document).ready(function() {
    $('#cycleForm').submit(function(event) {
        event.preventDefault();

        if ($.trim($("#inputScene").val()) === "" ) {
            alert('Please select a Scene.');
            return false;
        } else
        if ($.trim($("#inputCycles").val()) === "" || !$.isNumeric($("#inputCycles").val())) {
            alert('Please enter a numeric value for Scene Cycles.');
            return false;
        } else
        if ($.trim($("#inputEndDelay").val()) === "" || !$.isNumeric($("#inputEndDelay").val())) {
            alert('Please enter a numeric value for End Delay even if it is 0.');
            return false;
        }

        var formData = JSON.stringify($(this).serializeFormJSON());
        $.post("api/v1/show/{{.ShowID}}/cycle", formData, function(data) {
            if (data.Error != false) {
                alert("Error: " + data.Message);
            } else {
                populateContent();
                $("#addCycleModal").dialog("close");
            }
        }, "json");

        return false;
    });
});
</script>
{{end}}