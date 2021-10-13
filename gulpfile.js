const gulp = require("gulp");
const uglify = require("gulp-uglify");
const concat = require("gulp-concat");

gulp.task("js", function () {
  gulp
    .src("./public/*.js")
    .pipe(uglify())
    .pipe(concat("bundle.min.js"))
    .pipe(gulp.dest("./public"));
});
