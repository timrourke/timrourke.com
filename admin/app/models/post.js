import Model from 'ember-data/model';
import attr from 'ember-data/attr';
import { belongsTo } from 'ember-data/relationships';

export default Model.extend({
  createdAt:  attr('date'),
  updatedAt:  attr('date'),
  title:      attr('string'),
  excerpt:    attr('string'),
  content:    attr('string'),
  permalink:  attr('string'),

  user:       belongsTo('user'),
});