package caterpillar

const CTPLR_TMPL = `
{{define "MENUBAR_VIEW_BTNS"}}
		<a href="{{.EditURL}}" style="text-decoration: none;">
			<button type="button" class="btn btn-default btn-md">
				<span class="glyphicon glyphicon-pencil"></span>
				<div class="ctpl_tooltip">
	        <div class="ctpl_tooltip_arrow ctpl_tooltip_arrow_top"></div>
	        <div class="ctpl_tooltip_inner">編集</div>
	      </div>
			</button>
		</a>
		<a href="{{.PagesURL}}" style="text-decoration: none;">
			<button type="button" class="btn btn-default btn-md">
				<span class="glyphicon glyphicon-th-list"></span>
				<div class="ctpl_tooltip">
	        <div class="ctpl_tooltip_arrow ctpl_tooltip_arrow_top"></div>
	        <div class="ctpl_tooltip_inner">ページ一覧</div>
	      </div>
			</button>
		</a>
		<a href="{{.LogoutURL}}" style="text-decoration: none;">
			<button type="button" class="btn btn-default btn-md">
				<span class="glyphicon glyphicon-log-out"></span>
				<div class="ctpl_tooltip">
	        <div class="ctpl_tooltip_arrow ctpl_tooltip_arrow_top"></div>
	        <div class="ctpl_tooltip_inner">ログアウト</div>
	      </div>
			</button>
		</a>
{{end}}
{{define "MENUBAR_EDIT_BTNS"}}
		<button type="button" id="ctpl_edit_submit" class="btn btn-default btn-md">
			<span class="glyphicon glyphicon-save"></span>
			<div class="ctpl_tooltip">
        <div class="ctpl_tooltip_arrow ctpl_tooltip_arrow_top"></div>
        <div class="ctpl_tooltip_inner">保存</div>
      </div>
		</button>
		<a href="{{.ViewURL}}" style="text-decoration: none;" title="変更破棄">
			<button type="button" class="btn btn-default btn-md">
				<span class="glyphicon glyphicon-trash"></span>
				<div class="ctpl_tooltip">
	        <div class="ctpl_tooltip_arrow ctpl_tooltip_arrow_top"></div>
	        <div class="ctpl_tooltip_inner">変更破棄</div>
	      </div>
			</button>
		</a>
{{end}}
{{define "CATERPILLAR"}}
	<p id="ctpl_menu" style="position: fixed; right: 10px; bottom: 20px; clear: both; margin:0px;">
		{{if .Edit}}{{template "MENUBAR_EDIT_BTNS" .}}{{else}}{{template "MENUBAR_VIEW_BTNS" .}}{{end}}
	</p>
	<!-- TODO remove bootstrap dependency. It may cause conflict with users design. -->
	<link href="/caterpillar/static/bootstrap/css/bootstrap-theme.min.css" rel="stylesheet" type="text/css"></link>
	<link href="/caterpillar/static/bootstrap/css/bootstrap.css" rel="stylesheet" type="text/css"></link>
	<link href="/caterpillar/static/caterpillar.css" rel="stylesheet" type="text/css"></link>
	<script src="/caterpillar/static/tinymce/tinymce.min.js" type="text/javascript"></script>
	<script src="/caterpillar/static/tinymce_init.js" type="text/javascript"></script>
	<script type="text/javascript">
	var caterpillar = {};
	caterpillar.pageId = {{.PageID}};
	caterpillar.viewUrl = '{{.ViewURL}}';
	caterpillar.editUrl = '{{.EditURL}}';
	</script>
	<link href="/caterpillar/static/tinymce.css" rel="stylesheet" type="text/css"></link>
{{end}}
`
