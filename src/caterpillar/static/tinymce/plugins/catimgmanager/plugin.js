tinymce.PluginManager.add('catimgmanager', function(editor, url) {
	function FormBox() {
		editor.windowManager.open({
			title: '画像アップロード',
			url: url + '/dialog.html',
			width: 350,
			height: 240,
			buttons: [
				{text: '閉じる', onclick: 'close'}
			]
		});
	}

	editor.addButton('catimgmanager', {
		text: '画像アップロード',
		icon: 'image',
		onclick: FormBox
	});
	editor.addMenuItem('catimgmanager', {
		text: '画像アップロード',
		icon: 'image',
		context: 'file',
		onclick: FormBox
	});
});
