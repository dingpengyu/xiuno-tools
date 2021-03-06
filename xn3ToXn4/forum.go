package xn3ToXn4

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/skiy/golib"
	"log"
)

type forum struct {
	db3str,
	db4str dbstr
	fields forumFields
}

type forumFields struct {
	fid, name, rank, threads, todayposts, todaythreads, brief, accesson, orderby, icon, moduids, seo_title, seo_keywords string
}

func (this *forum) update() {
	if !lib.AutoUpdate(this.db4str.Auto, this.db4str.DBPre+"forum") {
		return
	}

	count, err := this.toUpdate()
	if err != nil {
		log.Fatalln("转换 " + this.db3str.DBPre + "forum 失败: " + err.Error())
	}

	fmt.Printf("转换 %sforum 表成功，共(%d)条数据\r\n", this.db3str.DBPre, count)
}

func (this *forum) toUpdate() (count int, err error) {
	xn3pre := this.db3str.DBPre
	xn4pre := this.db4str.DBPre

	fields := "fid,name,rank,threads,todayposts,todaythreads,brief,accesson,orderby,icon,moduids,seo_title,seo_keywords"
	qmark := this.db3str.FieldMakeQmark(fields, "?")
	xn3 := fmt.Sprintf("SELECT %s FROM %sforum", fields, xn3pre)
	xn4 := fmt.Sprintf("INSERT INTO %sforum (%s,announcement) VALUES (%s, '')", xn4pre, fields, qmark)

	data, err := xiuno3db.Query(xn3)
	if err != nil {
		log.Fatalln(xn3, err.Error())
	}
	defer data.Close()

	xn4Clear := "TRUNCATE `" + xn4pre + "forum`"
	_, err = xiuno4db.Exec(xn4Clear)
	if err != nil {
		log.Fatalf(":::清空 %sforum 表失败: "+err.Error(), xn4pre)
	}
	fmt.Printf("清空 %sforum 表成功\r\n", xn4pre)

	tx, err := xiuno4db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(xn4)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	fmt.Printf("正在升级 %sforum 表\r\n", xn4pre)

	var field forumFields
	for data.Next() {
		err = data.Scan(
			&field.fid,
			&field.name,
			&field.rank,
			&field.threads,
			&field.todayposts,
			&field.todaythreads,
			&field.brief,
			&field.accesson,
			&field.orderby,
			&field.icon,
			&field.moduids,
			&field.seo_title,
			&field.seo_keywords)

		_, err = stmt.Exec(
			&field.fid,
			&field.name,
			&field.rank,
			&field.threads,
			&field.todayposts,
			&field.todaythreads,
			&field.brief,
			&field.accesson,
			&field.orderby,
			&field.icon,
			&field.moduids,
			&field.seo_title,
			&field.seo_keywords)

		if err != nil {
			fmt.Printf("导入数据失败(%s) \r\n", err.Error())
		} else {
			count++
		}
	}

	if err = data.Err(); err != nil {
		log.Fatalln(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalln(err.Error())
	}

	return count, err
}
