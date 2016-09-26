import Ember from 'ember';

const { Route } = Ember;
const { Promise } = Ember.RSVP;

export default Route.extend({
  model(params) {
		let id = params.post_id;

		return new Promise((resolve, reject) => {
			this.store.query('post', { 
				filter: {
					permalink: id 
				},
				page: {
					limit: 1
				}
			}).then(posts => {
				// Attempt loading post by permalink first
				if (posts.get('firstObject')) {
					resolve(posts.get('firstObject')); 
				} else {
					// If that fails, try looking up by id
					this.store.findRecord('post', id).then(post => {
						this.transitionTo('post', post.get('permalink')); 
						resolve(post);
					});
				}
			}).catch(err => reject(err));
		});
  },

  setupController(controller, model) {
    controller.set('post', model);
  }
});
