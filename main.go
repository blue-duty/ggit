package main

import (
	"fmt"
	"github.com/emirpasic/gods/trees/binaryheap"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	commitgraph_fmt "github.com/go-git/go-git/v5/plumbing/format/commitgraph"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/object/commitgraph"
	"got/cmd"
	"io"
	"os"
	"path"
	"strings"
)

// Example how to resolve a revision into its commit counterpart
func main() {
	//CheckArgs("<path>", "<revision>", "<tree path>")

	//path := os.Args[1]
	//revision := os.Args[2]
	//treePath := os.Args[3]

	cmd.Execute()
	//r, err := exec.Command("git", "show", "--color", "--pretty=format:%b", "5c5b17f", "cmd/commit.go").Output()
	//if err != nil {
	//	panic(err)
	//}

	// 将结果绘制在新终端中并等待其退出
	//cmd := exec.Command("git", "show", "--color", "--pretty=format:%b", "5c5b17f", "cmd/commit.go")
	//cmd.Stdout = os.Stdout
	//cmd.Stdin = os.Stdin
	//cmd.Stderr = os.Stderr
	//err := cmd.Run()
	//if err != nil {
	//	return
	//}

	//cmd := exec.Command("ls", "-l")
	//cmd.Stdout = os.Stdout
	//cmd.Stdin = os.Stdin
	//cmd.Stderr = os.Stderr
	//err := cmd.Run()
	//if err != nil {
	//	return
	//}
	//
	//fmt.Println(cmd.Dir)
	//fmt.Println(cmd.Env)
	//fmt.Println(cmd.Process)
	//fmt.Println(cmd.ProcessState)
	//fmt.Println(cmd.String())
	//fmt.Println(cmd.CombinedOutput())
	//fmt.Println("push success")

	//fmt.Println("测试")

	//println(string(r))
	//println(""

	//path := "/home/duty/go/src/got"
	//
	//r, err := git.PlainOpen(path)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//// 获取工作目录状态
	//w, err := r.Worktree()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//status, err := w.Status()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	// 遍历修改过的文件
	//for file, s := range status {
	//if s.Worktree != git.Unmodified {
	//	filePath := filepath.Join(w.Filesystem.Root(), file)
	//	//f, err := os.Open(filePath)
	//	//if err != nil {
	//	//	fmt.Println(err)
	//	//	continue
	//	//}
	//	//defer f.Close()
	//
	//	// 获取文件内容
	//	//fileContents, err := getContents(f)
	//	//if err != nil {
	//	//	fmt.Println(err)
	//	//	continue
	//	//}
	//
	//	// 处理文件内容
	//	fmt.Println(file)
	//}

	//	if s.Worktree == git.Unmodified {
	//		// 判断文件是否存在
	//		//filePath := filepath.Join(w.Filesystem.Root(), file)
	//		if s.Staging == git.Deleted {
	//			fmt.Println("Deleted:", file)
	//			continue
	//		}
	//	}
	//	switch s.Worktree {
	//	case git.Unmodified:
	//		fmt.Println("Unmodified:", file)
	//	case git.Added:
	//		fmt.Println("Added:", file)
	//	case git.Deleted:
	//		fmt.Println("Deleted:", file)
	//	case git.Modified:
	//		fmt.Println("Modified:", file)
	//	case git.Renamed:
	//		fmt.Println("Renamed:", file)
	//	case git.Copied:
	//		fmt.Println("Copied:", file)
	//	case git.Untracked:
	//		fmt.Println("Untracked:", file)
	//	default:
	//		fmt.Println("Unknown:", file)
	//	}
	//}

	// 提交被删除的文件，不提交修改的文件
	//for file, s := range status {
	//	if s.Worktree == git.Unmodified {
	//		// 判断文件是否存在
	//		//filePath := filepath.Join(w.Filesystem.Root(), file)
	//		if s.Staging == git.Deleted {
	//			//fmt.Println("Deleted:", file)
	//			//continue
	//			// 提交被删除的文件
	//			_, err := w.Add(file)
	//			if err != nil {
	//				fmt.Println(err)
	//				return
	//			}
	//		}
	//	}

	//// 提交文件
	//_, err := w.Add(file)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//}

	//_, err = w.Add("main.go")
	//CheckIfError(err)
	//
	//// 提交
	//commit, err := w.Commit("commit", &git.CommitOptions{
	//	Author: &object.Signature{
	//		Name:  "John Doe",
	//		Email: "john@doe.org",
	//		When:  time.Now(),
	//	},
	//})
	//CheckIfError(err)
	//
	//obj, err := r.CommitObject(commit)
	//CheckIfError(err)
	//
	//fmt.Println(obj)

	// 遍历未提交的文件
	//var files []string
	//for file, s := range status {
	//	if s.Worktree != git.Unmodified {
	//		// 提交文件
	//		_, err := w.Add(file)
	//		if err != nil {
	//			fmt.Println(err)
	//			return
	//		}
	//		files = append(files, file)
	//	} else {
	//		// 提交被删除的文件
	//		if s.Staging == git.Deleted {
	//			// 提交文件
	//			_, err := w.Add(file)
	//			if err != nil {
	//				fmt.Println(err)
	//				return
	//			}
	//			files = append(files, file)
	//		}
	//	}
	//}
	//
	//// 提交
	//_, err = w.Commit("testcommit", &git.CommitOptions{})
	//CheckIfError(err)

	// 打印提交的文件
	//fmt.Println("Commit files:")
	//for _, file := range files {
	//	fmt.Println(file)
	//}
	//
	//fmt.Println("Done")

	// 获取最新引用
	//ref, err := r.Head()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//// 获取最新提交对象
	//commit, err := r.CommitObject(ref.Hash())
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//// 显示提交信息
	////fmt.Println(commit)
	//
	//// 显示文件更改
	//treeBefore, err := commit.Tree()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//parent, err := commit.Parent(0)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//treeAfter, err := parent.Tree()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//changes, err := object.DiffTreeWithOptions(context.TODO(), , treeBefore, &object.DiffTreeOptions{
	//	DetectRenames: true,
	//})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	////for _, change := range changes {
	////	fmt.Println(change.String())
	////}
	//
	//fmt.Println(changes.String())

	// We instantiate a new repository targeting the given path (the .git folder)
	//fs := osfs.New(path)
	//if _, err := fs.Stat(git.GitDirName); err == nil {
	//	fs, err = fs.Chroot(git.GitDirName)
	//	CheckIfError(err)
	//}

	//s := filesystem.NewStorageWithOptions(fs, cache.NewObjectLRUDefault(), filesystem.Options{KeepDescriptors: true})
	//r, err := git.Open(s, fs)
	//CheckIfError(err)
	//defer s.Close()

	// Resolve revision into a sha1 commit, only some revisions are resolved
	// look at the doc to get more details
	//Info("git rev-parse %s", revision)

	//h, err := r.ResolveRevision(plumbing.Revision(revision))
	//CheckIfError(err)

	//wt, err := r.Worktree()
	//CheckIfError(err)
	//
	//status, err := wt.Status()
	//CheckIfError(err)
	//
	////fmt.Println(status.String())
	//
	//changedFiles := make([]string, 0)
	//for file, fileStatus := range status {
	//	if fileStatus.Worktree == git.Modified {
	//		changedFiles = append(changedFiles, file)
	//	}
	//}
	//
	//fmt.Println(changedFiles)

	//for path, file := range status {
	//	//fmt.Println(path, file.Staging, file.Worktree)
	//	fmt.Printl
	//}

	//ref, err := r.Head()
	//CheckIfError(err)

	// ... retrieves the commit history
	//commit, err := r.Log(&git.LogOptions{From: ref.Hash()})
	//CheckIfError(err)
	//h := ref.Hash()

	//commit, err := r.CommitObject(ref.Hash())
	//CheckIfError(err)
	//
	//tree, err := commit.Tree()
	//CheckIfError(err)
	//if treePath != "" {
	//	tree, err = tree.Tree(treePath)
	//	CheckIfError(err)
	//}

	//var paths []string
	//for _, entry := range tree.Entries {
	//	paths = append(paths, entry.Name)
	//}
	//
	//commitNodeIndex, file := getCommitNodeIndex(r, fs)
	//if file != nil {
	//	defer file.Close()
	//}

	//commitNode, err := commitNodeIndex.Get(ref.Hash())
	//CheckIfError(err)

	//revs, err := getLastCommitForPaths(commitNode, "", paths)
	//CheckIfError(err)
	//for path, rev := range revs {
	//	// Print one line per file (name hash message)
	//	hash := rev.Hash.String()
	//	line := strings.Split(rev.Message, " \n")
	//	//fmt.Println(path, hash[:7], line[0])
	//}
}

func getCommitNodeIndex(r *git.Repository, fs billy.Filesystem) (commitgraph.CommitNodeIndex, io.ReadCloser) {
	file, err := fs.Open(path.Join("objects", "info", "commit-graph"))
	if err == nil {
		index, err := commitgraph_fmt.OpenFileIndex(file)
		if err == nil {
			return commitgraph.NewGraphCommitNodeIndex(index, r.Storer), file
		}
		file.Close()
	}

	return commitgraph.NewObjectCommitNodeIndex(r.Storer), nil
}

type commitAndPaths struct {
	commit commitgraph.CommitNode
	// Paths that are still on the branch represented by commit
	paths []string
	// Set of hashes for the paths
	hashes map[string]plumbing.Hash
}

func getCommitTree(c commitgraph.CommitNode, treePath string) (*object.Tree, error) {
	tree, err := c.Tree()
	if err != nil {
		return nil, err
	}

	// Optimize deep traversals by focusing only on the specific tree
	if treePath != "" {
		tree, err = tree.Tree(treePath)
		if err != nil {
			return nil, err
		}
	}

	return tree, nil
}

func getFullPath(treePath, path string) string {
	if treePath != "" {
		if path != "" {
			return treePath + "/" + path
		}
		return treePath
	}
	return path
}

func getFileHashes(c commitgraph.CommitNode, treePath string, paths []string) (map[string]plumbing.Hash, error) {
	tree, err := getCommitTree(c, treePath)
	if err == object.ErrDirectoryNotFound {
		// The whole tree didn't exist, so return empty map
		return make(map[string]plumbing.Hash), nil
	}
	if err != nil {
		return nil, err
	}

	hashes := make(map[string]plumbing.Hash)
	for _, path := range paths {
		if path != "" {
			entry, err := tree.FindEntry(path)
			if err == nil {
				hashes[path] = entry.Hash
			}
		} else {
			hashes[path] = tree.Hash
		}
	}

	return hashes, nil
}

func getLastCommitForPaths(c commitgraph.CommitNode, treePath string, paths []string) (map[string]*object.Commit, error) {
	// We do a tree traversal with nodes sorted by commit time
	heap := binaryheap.NewWith(func(a, b interface{}) int {
		if a.(*commitAndPaths).commit.CommitTime().Before(b.(*commitAndPaths).commit.CommitTime()) {
			return 1
		}
		return -1
	})

	resultNodes := make(map[string]commitgraph.CommitNode)
	initialHashes, err := getFileHashes(c, treePath, paths)
	if err != nil {
		return nil, err
	}

	// Start search from the root commit and with full set of paths
	heap.Push(&commitAndPaths{c, paths, initialHashes})

	for {
		cIn, ok := heap.Pop()
		if !ok {
			break
		}
		current := cIn.(*commitAndPaths)

		// Load the parent commits for the one we are currently examining
		numParents := current.commit.NumParents()
		var parents []commitgraph.CommitNode
		for i := 0; i < numParents; i++ {
			parent, err := current.commit.ParentNode(i)
			if err != nil {
				break
			}
			parents = append(parents, parent)
		}

		// Examine the current commit and set of interesting paths
		pathUnchanged := make([]bool, len(current.paths))
		parentHashes := make([]map[string]plumbing.Hash, len(parents))
		for j, parent := range parents {
			parentHashes[j], err = getFileHashes(parent, treePath, current.paths)
			if err != nil {
				break
			}

			for i, path := range current.paths {
				if parentHashes[j][path] == current.hashes[path] {
					pathUnchanged[i] = true
				}
			}
		}

		var remainingPaths []string
		for i, path := range current.paths {
			// The results could already contain some newer change for the same path,
			// so don't override that and bail out on the file early.
			if resultNodes[path] == nil {
				if pathUnchanged[i] {
					// The path existed with the same hash in at least one parent so it could
					// not have been changed in this commit directly.
					remainingPaths = append(remainingPaths, path)
				} else {
					// There are few possible cases how can we get here:
					// - The path didn't exist in any parent, so it must have been created by
					//   this commit.
					// - The path did exist in the parent commit, but the hash of the file has
					//   changed.
					// - We are looking at a merge commit and the hash of the file doesn't
					//   match any of the hashes being merged. This is more common for directories,
					//   but it can also happen if a file is changed through conflict resolution.
					resultNodes[path] = current.commit
				}
			}
		}

		if len(remainingPaths) > 0 {
			// Add the parent nodes along with remaining paths to the heap for further
			// processing.
			for j, parent := range parents {
				// Combine remainingPath with paths available on the parent branch
				// and make union of them
				remainingPathsForParent := make([]string, 0, len(remainingPaths))
				newRemainingPaths := make([]string, 0, len(remainingPaths))
				for _, path := range remainingPaths {
					if parentHashes[j][path] == current.hashes[path] {
						remainingPathsForParent = append(remainingPathsForParent, path)
					} else {
						newRemainingPaths = append(newRemainingPaths, path)
					}
				}

				if remainingPathsForParent != nil {
					heap.Push(&commitAndPaths{parent, remainingPathsForParent, parentHashes[j]})
				}

				if len(newRemainingPaths) == 0 {
					break
				} else {
					remainingPaths = newRemainingPaths
				}
			}
		}
	}

	// Post-processing
	result := make(map[string]*object.Commit)
	for path, commitNode := range resultNodes {
		var err error
		result[path], err = commitNode.Commit()
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		Warning("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		os.Exit(1)
	}
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}
