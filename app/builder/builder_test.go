package builder_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/Unknwon/com"
	"github.com/go-xiaohei/pugo-static/app/builder"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	opt = &builder.BuildOption{
		SrcDir:   "../../source",
		TplDir:   "../../template",
		MediaDir: "../../source/media",
		Version:  "0.1.1",
		VerDate:  "2015-12-22",
		Theme:    "default",
	}
	target          = "../../dest"
	b               = builder.New(opt)
	shoudlExistDirs = []string{
		"tags",
		"posts",
		"static",
	}
	shouldExistFiles = []string{
		"index.html",
		"archive.html",
		"about.html",
		"feed.xml",
		"sitemap.xml",
		"favicon.ico",
		"tags/pugo.html",
		"posts/1.html",
	}
)

func TestBuilderBuild(t *testing.T) {
	Convey("Build Process", t, func() {
		b.Build(target)
		So(b.Error, ShouldBeNil)
		So(b.Context(), ShouldNotBeNil)
		So(b.Context().Error, ShouldBeNil)

		// check dirs and files
		Convey("Check Built And Files", func() {
			var flag = true
			t := path.Join(target, b.Context().Meta.Base)
			for _, dir := range shoudlExistDirs {
				if flag = flag && com.IsDir(path.Join(t, dir)); !flag {
					break
				}
			}
			So(flag, ShouldBeTrue)

			for _, f := range shouldExistFiles {
				if flag = flag && com.IsFile(path.Join(t, f)); !flag {
					break
				}
			}
			So(flag, ShouldBeTrue)
		})
	})
}

func TestBuilderErrors(t *testing.T) {
	Convey("Build Errors", t, func() {
		opt.SrcDir = "../../xxxx"
		b2 := builder.New(opt)
		So(b2.Error, ShouldEqual, builder.ErrSrcDirMissing)

		opt.SrcDir = "../../source"
		opt.TplDir = "../../xxx"
		b2 = builder.New(opt)
		So(b2.Error, ShouldEqual, builder.ErrTplDirMissing)

	})

	Convey("Build Fail", t, func() {
		opt.TplDir = "../../template"
		opt.SrcDir = "./testdata"

		b := builder.New(opt)
		So(b.Error, ShouldBeNil)

		b.Build("testdata_dest")
		So(b.Context().Error, ShouldNotBeNil)

		removeDirectory("testdata_dest")
	})
}

/*
// it seems always fail in drone.io
func TestBuildWatch(t *testing.T) {
	Convey("Build Watch", t, func() {
		opt.TplDir = "../template"
		opt.SrcDir = "./testdata"
		b := builder.New(opt)
		So(b.Error, ShouldBeNil)
		So(b.Context(), ShouldBeNil)

		b.Watch("testdata_dest")
		file := path.Join(opt.SrcDir, "test.md")
		ioutil.WriteFile(file, []byte("```ini"), os.ModePerm)

		time.Sleep(5 * time.Second)
		So(b.Context(), ShouldNotBeNil)
		So(b.Context().Error, ShouldNotBeNil)
	})
}*/

// remove all sub dirs and files in directory
func removeDirectory(dir string) error {
	if !com.IsDir(dir) {
		return nil
	}
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, d := range dirs {
		if d.IsDir() {
			if err = removeDirectory(path.Join(dir, d.Name())); err != nil {
				return err
			}
		}
	}
	return os.RemoveAll(dir)
}