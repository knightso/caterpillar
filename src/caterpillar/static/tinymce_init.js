(function() {
	var editedIDs = {};	// 編集したエディタのID一覧。

	tinymce.init({
		inline: true,
		language: 'ja',
		selector: 'div.ctpl_block',
		plugins: [
		  'advlist autolink lists link image charmap print preview hr anchor',
		  'searchreplace visualblocks code fullscreen',
		  'insertdatetime media table contextmenu paste textcolor colorpicker'
		],
		toolbar: 'undo redo | styleselect | bold italic forecolor backcolor | alignleft aligncenter alignright alignjustify | bullist numlist outdent indent | anchor link image',
		setup: function(editor) {
			// 編集したエディタを捕捉する。
			editor.on('blur', function(e) {
				if(editor.id in editedIDs) {
					return;
				}
				editedIDs[editor.id] = null;
			});
		},
		file_picker_callback: function(callback, value, meta) {
			var editor = tinymce.activeEditor;
			var width = document.documentElement.clientWidth * 0.8;
			var height = document.documentElement.clientHeight - 90;
			var fileWin;

			if(meta.filetype == 'image') {
				// imageプラグインからの呼び出し。
				fileWin = editor.windowManager.open({
					title: '画像選択',
					url: '/caterpillar/static/tinymce/plugins/catimgmanager/filebrowser/filebrowser.html',
					width: width,
					height: height,
					resizable: true,
					scrollbars: true,
					buttons: [
						{text: 'キャンセル', onclick: 'close'}
					]
				});
				editor.windowManager.setParams({
					requestUrl: '/caterpillar/filemanager/files',
					callback: callback
				});
			} else if(meta.filetype == 'file') {
				// linkプラグインからの呼び出し。
				fileWin = editor.windowManager.open({
					title: 'ページ選択',
					url: '/caterpillar/static/tinymce/plugins/catpageselector/dialog.html',
					width: width,
					height: height,
					resizable: true,
					scrollbars: true,
					buttons: [
						{text: 'キャンセル', onclick: 'close'}
					]
				});
				editor.windowManager.setParams({
					requestUrl: '/caterpillar/api/pages',
					callback: callback
				});
			}

			// ウィンドウのリサイズイベントを登録します。
			var eventQueue;
			window.addEventListener('resize', function() {
				if(eventQueue !== false) {
					clearTimeout(eventQueue);
				}
				eventQueue = setTimeout(function() {
					// ウィンドウサイズの調整。
					var width = parseInt(document.documentElement.clientWidth * 0.8);
					var height = parseInt(document.documentElement.clientHeight);
					fileWin.width(width);
					fileWin.height(height);
					fileWin.resizeTo(width, height);
					// ウィンドウの位置調整。
					var width = document.documentElement.clientWidth;
					var fileWinWidth = fileWin.width();
					fileWin.moveTo((width-fileWinWidth)/2, 0);
				}, 300);
			});
		}
	});

	// エディタの編集したコンテンツをサーバーに送信します。
	// W3C DOCイベントモデルサポートブラウザのみ対応 FireFox, Chrome, Safari, Opera, IE9 ～
	window.addEventListener('load', function() {
		var submitButton = document.getElementById('ctpl_edit_submit');
		submitButton.addEventListener('click', function(evt) {
			var blocks = {
			};

			var edited = false;
			for(var id in editedIDs) {
				var editor = tinymce.EditorManager.get(id);
				blocks[id] = editor.getContent();
				edited = true;
			}
			if (!edited) {
				// TODO disable save button until modified.
				window.alert('修正されていません。');
				return;
			}

			tinymce.util.XHR.send({
				url: '/caterpillar/api/blocks/' + caterpillar.pageId,
				type: 'PUT',
				data: tinymce.util.JSON.serialize(blocks),
				success: function() {
					window.alert('保存しました。');
					document.location=caterpillar.viewUrl;
				},
				error: function() {
					window.alert('予期せぬエラーが発生しました。。');
				}
			});
		}, false);
	});

})();
