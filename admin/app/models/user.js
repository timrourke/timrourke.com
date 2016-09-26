import Model from 'ember-data/model';
import attr from 'ember-data/attr';
import { hasMany } from 'ember-data/relationships';

export default Model.extend({
  createdAt:  attr('date'),
  updatedAt:  attr('date'),
  email:      attr('string'),
  username:   attr('string'),

  posts:      hasMany('post'),
});