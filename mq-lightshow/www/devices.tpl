{{define "content"}}
<h1>Devices</h1>
<table class="table">
  <thead>
    <tr>
      <th scope="col">Name</th>
      <th scope="col">Topic</th>
      <th scope="col">Type</th>
      <th scope="col"><a href="devices-add" class="btn btn-sm btn-primary">Add</a></th>
    </tr>
  </thead>
  <tbody>
{{ range .Devices }}
    <tr>
      <td>{{.Name}}</td>
      <td>{{.Topic}}</td>
      <td>{{.Type.Name}}</td>
      <td>
        <a href="devices-delete?deviceID={{.ID}}" class="btn btn-sm btn-danger" title="Delete Device" onclick="return confirm('Are you sure you want to delete this device? Note that this may cause problems if the device is used in any scenes!')">
          <div class="icon-button-delete">&nbsp;</div>
        </a>
        <a href="devices-edit?deviceID={{.ID}}" class="btn btn-sm btn-primary">Edit</a>
      </td>
    </tr>
{{end}}
  </tbody>
</table>
<script>
$(document).ready(function() {
});
</script>
{{end}}
