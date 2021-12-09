{{define "base"}}<!DOCTYPE html>
<html lang="en">
  <head>
    <title>{{.PageInfo.Title}}</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta name="google" content="notranslate">
    <link rel="stylesheet" href="css/bootstrap.min.css">
    <link rel="stylesheet" href="css/jquery-ui.min.css">
    <link rel="stylesheet" href="css/style.css">
    <script src="js/jquery-3.4.1.min.js"></script>
    <script src="js/bootstrap.min.js"></script>
    <script src="js/jquery-ui.min.js"></script>
  </head>
  <body>
    <nav class="navbar navbar-expand-md navbar-dark bg-dark">
        <a class="navbar-brand" href="shows">Light Shows</a>
        <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbars" aria-controls="navbars" aria-expanded="false" aria-label="Toggle navigation">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbars">
          <ul class="navbar-nav mr-auto">
            <li class="nav-item{{if .PageInfo.ScenesLinkEnabled}} active{{end}}">
              <a class="nav-link" href="scenes">Scenes {{if .PageInfo.ScenesLinkEnabled}}<span class="sr-only">(current)</span>{{end}}</a>
            </li>
            <li class="nav-item{{if .PageInfo.DevicesLinkEnabled}} active{{end}}">
              <a class="nav-link" href="devices">Devices {{if .PageInfo.DevicesLinkEnabled}}<span class="sr-only">(current)</span>{{end}}</a>
            </li>
            <li class="nav-item{{if .PageInfo.MQTTLinkEnabled}} active{{end}}">
              <a class="nav-link" href="mqtt">MQTT {{if .PageInfo.MQTTLinkEnabled}}<span class="sr-only">(current)</span>{{end}}</a>
            </li>
          </ul>
        </div>
    </nav>
    <main class="container">
        <div class="container" id="mainContainer">
{{template "content" .}}
<div id="focusSelector" style="display:none"></div>
        <br>
        </div>
    </main>
<script>
(function ($) {
    $.fn.serializeFormJSON = function () {
        var o = {};
        var a = this.serializeArray();
        $.each(a, function () {
            if (o[this.name]) {
                if (!o[this.name].push) {
                    o[this.name] = [o[this.name]];
                }
                o[this.name].push(this.value || '');
            } else {
                o[this.name] = this.value || '';
            }
        });
        return o;
    };
})(jQuery);

function makeDialog(selector) {
  if ($(selector).hasClass('ui-dialog-content')) {
    $(selector).dialog('destroy');
  }

  $(selector).dialog({
    autoOpen: false,
    position: { my: "center top", at: "center top", of:  this },
    classes: {
      "ui-dialog": "ui-corner-all",
      "ui-dialog-titlebar": "ui-corner-all bg-primary text-white",
    },
    close: function( event, ui ) {
      $(selector).empty();
    }
  });
}

function focusPage() {
  selector = $('#focusSelector').text();
  if (selector != "") {
    offset = $(selector).offset();
    if (offset) {
      pos = offset.top - 250
      $('html, body').animate({ scrollTop: pos }, 'fast');
    }
  }
}

</script>
  </body>
</html>{{end}}