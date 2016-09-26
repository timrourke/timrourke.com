import Ember from 'ember';

const { Component, run } = Ember;

export default Component.extend({
	classNames: ['quill-js'],

	content: "",
	editor: null,

	didInsertElement() {
		this._super(...arguments);
		
		let editor = new Quill('#editor', {
			modules: { toolbar: '#toolbar' },
			theme: 'snow'
		});

		editor.setText(this.get('content'));

		editor.on('text-change', () => {
			run(() => {
				this.set('content', editor.getText());
			});
		});

		this.set('editor', editor);
	}
});
