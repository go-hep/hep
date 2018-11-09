// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

const page = `<html>
<head>
    <title>go-hep/groot file inspector</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css" />
	<link rel="stylesheet" href="https://www.w3schools.com/w3css/3/w3.css">
	<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.1.1/jquery.min.js"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/jqueryui/1.12.1/jquery-ui.min.js"></script>
	<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/jstree/3.3.7/themes/default/style.min.css" />
	<script src="https://cdnjs.cloudflare.com/ajax/libs/jstree/3.3.7/jstree.min.js"></script>
	<style>
	input[type=file] {
		display: none;
	}
	input[type=submit] {
		background-color: #F44336;
		padding:5px 15px;
		border:0 none;
		cursor:pointer;
		-webkit-border-radius: 5px;
		border-radius: 5px;
	}
	.flex-container {
		display: -webkit-flex;
		display: flex;
	}
	.flex-item {
		margin: 5px;
	}
	.groot-file-upload {
		color: white;
		background-color: #0091EA;
		padding:5px 15px;
		border:0 none;
		cursor:pointer;
		-webkit-border-radius: 5px;
	}
	</style>
</head>
<body>

<!-- Sidebar -->
<div id="groot-sidebar" class="w3-sidebar w3-bar-block w3-card-4 w3-light-grey" style="width:25%">
	<div class="w3-bar-item w3-card-2 w3-black">
		<h2>go-hep/groot ROOT file inspector</h2>
	</div>
	<div class="w3-bar-item">

	{{if .Local}}
	<div>
		File: <input id="groot-open-form-input" type="text" name="uri" value placeholder="URI to local or remote file">
		<label for="groot-open-button" class="groot-file-upload" style="font-size:16px" onclick="openROOTFile()">
		<i class="fa fa-folder-open" aria-hidden="true" style="font-size:16px"></i> Open
		</label>
		<input id="groot-open-button" type="hidden" value="Open" onclick="openROOTFile()">
	</div>
	<br>
	{{- end}}

	<form id="groot-upload-form" enctype="multipart/form-data" action="/root-file-upload" method="post">
		<label for="groot-file-upload" class="groot-file-upload" style="font-size:16px">
		<i class="fa fa-cloud-upload" aria-hidden="true" style="font-size:16px"></i> Upload
		</label>
		<input id="groot-file-upload" type="file" name="groot-file"/>
		<input type="hidden" name="token" value="{{.Token}}"/>
		<input type="hidden" value="upload" />
	</form>

	</div>
	<div id="groot-file-tree" class="w3-bar-item">
	</div>
</div>

<!-- Page Content -->
<div style="margin-left:25%; height:100%" class="w3-grey" id="groot-container">
<div class="w3-container w3-content w3-cell w3-cell-middle w3-cell-row w3-center w3-justify w3-grey" style="width:100%" id="groot-display">
	</div>
</div>

<script type="text/javascript">
	document.getElementById("groot-file-upload").onchange = function() {
		var data = new FormData($("#groot-upload-form")[0]);
		var dst = data.get("groot-file").name;
		dst = "upload://"+dst.substring(dst.lastIndexOf('/')+1);
		data.append("groot-dst", dst);
		$.ajax({
			url: "/root-file-upload",
			method: "POST",
			data: data,
			processData: false,
			contentType: false,
			success: displayFileTree,
			error: function(er){
				alert("upload failed: "+er);
			}
		});
	}
{{if .Local}}
	function openROOTFile() {
		var uri = $("#groot-open-form-input").val();
		$("#groot-open-form-input").val("");
		var data = new FormData();
		data.append("uri", uri);
		$.ajax({
			url: "/root-file-open",
			method: "POST",
			data: data,
			processData: false,
			contentType: false,
			success: displayFileTree,
			error: function(e){
				alert("open failed: "+e);
			}
		});
	}
	$('#groot-open-form-input').keypress(function(event) {
		if (event.keyCode == 13) {
			openROOTFile();
		}
	});

/*
	$('#groot-test-form').dialog({
        modal: true,
		show: false,
        buttons: {
            'Open-1': function () {
                var name = $('input[name="uri"]').val();
                $.post({
					url: "/file-open",
					method: "POST",
					data: {"uri": name},
					processData: false,
					contentType: "application/json",
					dataType: "json",
					success: displayFileTree,
					error: function(er){
						alert("open failed: "+JSON.stringify(er));
					}
				})
                $(this).dialog('close');
            },
			'Cancel': function () {
                $(this).dialog('close');
            }
        }
    });
*/
{{- end}}
	$(function () {
		$('#groot-file-tree').jstree();
		$("#groot-file-tree").on("select_node.jstree",
			function(evt, data){
				data.instance.toggle_node(data.node);
				if (data.node.a_attr.plot) {
					data.instance.deselect_node(data.node);
					$.post({
						type: 'POST',
						url: data.node.a_attr.href,
						data: data.node.a_attr.cmd,
						success: plotCallback,
						contentType: "application/json",
						dataType: 'json',
					});
				//	$.get(data.node.a_attr.href, plotCallback);
				}
			}
		);
		$.ajax({
			url: "/refresh",
			method: "GET",
			processData: false,
			contentType: false,
			success: displayFileTree,
			error: function(er){
				alert("refresh failed: "+er);
			}
		});
	});

	function displayFileTree(data) {
		$('#groot-file-tree').jstree(true).settings.core.data = JSON.parse(data);
		$("#groot-file-tree").jstree(true).refresh();
	};

	function plotCallback(data, status) {
		var img = data;
		var node = $("<div></div>");
		node.addClass("w3-panel w3-white w3-card-2 w3-display-container w3-content w3-center");
		node.css("width","100%");
		node.html(
			""
			+atob(img.data)
			+"<span onclick=\"this.parentElement.style.display='none'; updateHeight();\" class=\"w3-button w3-display-topright w3-hover-red w3-tiny\">X</span>"
		);
		$("#groot-display").prepend(node);
		updateHeight();
	};

	function updateHeight() {
		var hmenu = $("#groot-sidebar").height();
		var hcont = $("#groot-container").height();
		var hdisp = $("#groot-display").height();
		if (hdisp > hcont) {
			$("#groot-container").height(hdisp);
		}
		if (hdisp < hmenu && hcont > hmenu) {
			$("#groot-container").height(hmenu);
		}
	};
</script>
</body>
</html>
`
