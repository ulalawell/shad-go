//go:build !solution

package hogwarts

func GetCourseList(prereqs map[string][]string) []string {
	adjList := map[string][]string{}
	inpDegree := map[string]int{}
	visited := map[string]struct{}{}
	courseList := []string{}

	for course, prerequisites := range prereqs {
		for _, prerequisite := range prerequisites {
			adjList[prerequisite] = append(adjList[prerequisite], course)

			if _, ok := inpDegree[prerequisite]; !ok {
				inpDegree[prerequisite] = 0
			}

			inpDegree[course]++
		}
	}

	st := []string{}
	for course, degree := range inpDegree {
		if degree == 0 {
			st = append(st, course)
			visited[course] = struct{}{}
		}
	}

	for len(st) > 0 {
		fromCourse := st[len(st)-1]
		st = st[:len(st)-1]

		if inpDegree[fromCourse] > 0 {
			continue
		}

		visited[fromCourse] = struct{}{}
		courseList = append(courseList, fromCourse)

		for _, toCourse := range adjList[fromCourse] {
			inpDegree[toCourse]--
			st = append(st, toCourse)
		}
	}

	if len(visited) < len(inpDegree) {
		panic("Cycle found!")
	}

	return courseList
}
