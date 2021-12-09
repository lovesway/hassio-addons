{{define "content"}}
<form id="groupConfigureForm">
    <div class="form-group">
        <label for="inputDelay">Delay</label>
        <input type="text" class="form-control" id="inputDelay" aria-describedby="inputDelayHelp" name="Delay" value="{{.Group.Delay}}">
        <small id="inputDelayHelp" class="form-text text-muted">How much time to wait before executing the next group.</small>
    </div>
    <div class="form-group">
        <label for="inputGlobalDelay">Use Global Delay</label>
        <select class="custom-select" class="form-control" id="inputGlobalDelay" aria-describedby="inputGlobalDelayHelp" name="GlobalDelay">
            <option value="true"{{if (eq .Group.GlobalDelay true)}} selected{{end}}>true</option>
            <option value="false"{{if (eq .Group.GlobalDelay false)}} selected{{end}}>false</option>
        </select>
        <small id="inputGlobalDelayHelp" class="form-text text-muted">Use the Global Delay set in the Show configuration (only works when running in show context).</small>
    </div>
    <div class="form-group">
        <label for="inputName">Order</label>
        <input type="text" class="form-control" id="inputOrder" name="Order" value="{{.Group.Order}}">
    </div>
    <button type="submit" class="btn btn-primary">Submit</button>
</form>
<script>
$(document).ready(function() {
    $('#groupConfigureForm').submit(function(event) {
        event.preventDefault();

        if (!$.isNumeric($("#inputDelay").val())) {
            alert('Please enter a float value for the Delay time in seconds (5 seconds is 5.0).');
            return false;
        } else
        if ($.trim($("#inputOrder").val()) === "" || !$.isNumeric($("#inputOrder").val())) {
            alert('Please enter a numeric value for the order.');
            return false;
        }

        var formData = JSON.stringify($(this).serializeFormJSON());

        $.post("api/v1/scene/{{.SceneID}}/group/{{.Group.ID}}/configure", formData, function(data) {
            if (data.Error != false) {
                alert("Error: " + data.Message);
            } else {
                $('#focusSelector').text('#groupRow{{.Group.ID}}')
                populateContent();
                $("#configureGroupModal").dialog("close");
            }
        }, "json");

        return false;
    });    
});
</script>
{{end}}