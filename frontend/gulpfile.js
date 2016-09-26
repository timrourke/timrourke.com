'use strict';
 
var browserSync = require('browser-sync').create();
var gulp = require('gulp');
var notify = require('gulp-notify');
var rename = require('gulp-rename');
var sass = require('gulp-sass');
var sourcemaps = require('gulp-sourcemaps');
 
gulp.task('browser-sync', function() {
	browserSync.init({
		proxy: "localhost:8000",
	});
});

gulp.task('sass', function () {
 	return gulp.src('./sass/**/*.scss')
 		.pipe(sourcemaps.init())
		.pipe(sass({outputStyle: 'compressed'}).on('error', sass.logError))
		.pipe(sourcemaps.write())
		.pipe(gulp.dest('./html/css'))
		.pipe(rename('base.scss'))
		.pipe(gulp.dest('./../admin/app/styles'))
		.pipe(browserSync.stream())
		.pipe(notify("Compiled Sass."));
});
 
gulp.task('watch', function () {
	gulp.watch('./sass/**/*.scss', ['sass']);
	gulp.watch('./html/*.html').on('change', browserSync.reload);
});

gulp.task('default', ['sass', 'browser-sync', 'watch']);
