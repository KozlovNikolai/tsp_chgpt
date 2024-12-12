use std::cmp::Ordering;
use std::collections::BinaryHeap;
use std::time::Instant;

#[derive(Clone, Eq, PartialEq)]
struct Node {
    level: usize,
    path: Vec<usize>,
    bound: i32,
    current_cost: i32,
}

impl Ord for Node {
    fn cmp(&self, other: &Self) -> Ordering {
        // Reverse the order to create a min-heap
        other.bound.cmp(&self.bound)
    }
}

impl PartialOrd for Node {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

fn reduce_matrix(graph: &Vec<Vec<i32>>) -> (Vec<Vec<i32>>, i32) {
    let n = graph.len();
    let mut reduced = graph.clone();
    let mut lower_bound = 0;

    // Reduce rows
    for row in &mut reduced {
        let min_value = *row.iter().filter(|&&x| x != i32::MAX).min().unwrap_or(&i32::MAX);
        if min_value != i32::MAX {
            lower_bound += min_value;
            for value in row.iter_mut() {
                if *value != i32::MAX {
                    *value -= min_value;
                }
            }
        }
    }

    // Reduce columns
    for j in 0..n {
        let min_value = (0..n)
            .filter_map(|i| if reduced[i][j] != i32::MAX { Some(reduced[i][j]) } else { None })
            .min()
            .unwrap_or(i32::MAX);
        if min_value != i32::MAX {
            lower_bound += min_value;
            for i in 0..n {
                if reduced[i][j] != i32::MAX {
                    reduced[i][j] -= min_value;
                }
            }
        }
    }

    (reduced, lower_bound)
}

fn calculate_bound(graph: &Vec<Vec<i32>>, path: &Vec<usize>) -> i32 {
    let n = graph.len();
    let mut reduced = graph.clone();
    let mut lower_bound = 0;

    for i in 0..path.len() - 1 {
        reduced[path[i]][path[i + 1]] = i32::MAX;
        for j in 0..n {
            reduced[path[i]][j] = i32::MAX;
            reduced[j][path[i + 1]] = i32::MAX;
        }
    }

    let (_, additional_cost) = reduce_matrix(&reduced);
    lower_bound + additional_cost
}

fn contains(slice: &Vec<usize>, value: usize) -> bool {
    slice.contains(&value)
}

fn branch_and_bound(graph: Vec<Vec<i32>>, max_iterations: usize) -> (i32, Vec<usize>) {
    let n = graph.len();
    let (_, initial_lower_bound) = reduce_matrix(&graph);

    let mut pq = BinaryHeap::new();
    let initial_node = Node {
        level: 0,
        path: vec![0],
        bound: initial_lower_bound,
        current_cost: 0,
    };
    pq.push(initial_node);

    let mut best_cost = i32::MAX;
    let mut best_path = Vec::new();
    let mut iteration_count = 0;

    while let Some(current_node) = pq.pop() {
        if iteration_count >= max_iterations {
            println!("Maximum iterations reached. Terminating early.");
            break;
        }

        iteration_count += 1;

        if current_node.bound >= best_cost {
            continue;
        }

        if current_node.level == n - 1 {
            let final_cost = current_node.current_cost;
            if final_cost < best_cost {
                best_cost = final_cost;
                best_path = current_node.path.clone();
            }
            continue;
        }

        let current_city = *current_node.path.last().unwrap();
        for next_city in 0..n {
            if contains(&current_node.path, next_city) {
                continue;
            }

            let mut next_path = current_node.path.clone();
            next_path.push(next_city);
            let next_cost = current_node.current_cost + graph[current_city][next_city];
            let next_bound = calculate_bound(&graph, &next_path);

            if next_bound < best_cost {
                let next_node = Node {
                    level: current_node.level + 1,
                    path: next_path,
                    bound: next_bound,
                    current_cost: next_cost,
                };
                pq.push(next_node);
            }
        }
    }

    println!("Iteration count: {}", iteration_count);
    (best_cost, best_path)
}

fn calculate_distance(graph: &Vec<Vec<i32>>, path: &Vec<usize>) -> i32 {
    path.windows(2).map(|w| graph[w[0]][w[1]]).sum()
}

fn main() {
    const INF: i32 = i32::MAX;
    let graph = vec![
        vec![INF, 75987, 76499, 76503, 76657, 76726, 76554, 76917, 77002, 77221, 77562, 77770, 77774, 78402, 78439],
        vec![75627, INF, 909, 913, 1067, 1136, 964, 1328, 1412, 1631, 1972, 2181, 2184, 2813, 2849],
        vec![76194, 892, INF, 1141, 1296, 1364, 1192, 1556, 1641, 1859, 2201, 2409, 2412, 3041, 3077],
        vec![76204, 902, 1080, INF, 299, 368, 550, 913, 998, 1217, 1558, 1767, 1770, 2399, 2435],
        vec![76359, 1057, 1234, 299, INF, 173, 704, 1068, 1153, 1371, 1713, 1921, 1925, 2553, 2589],
        vec![76427, 1125, 1303, 368, 154, INF, 773, 1137, 1221, 1440, 1781, 1990, 1993, 2622, 2658],
        vec![76209, 907, 1084, 623, 777, 846, INF, 517, 602, 821, 1162, 1371, 1374, 2003, 2039],
        vec![76572, 1270, 1448, 987, 1141, 1210, 517, INF, 271, 490, 666, 874, 877, 1506, 1542],
        vec![76657, 1355, 1533, 1072, 1226, 1295, 602, 271, INF, 218, 560, 768, 771, 1400, 1436],
        vec![76876, 1574, 1751, 1290, 1444, 1513, 821, 490, 218, INF, 341, 549, 553, 1181, 1218],
        vec![77217, 1915, 2093, 1632, 1786, 1855, 1162, 666, 560, 341, INF, 208, 211, 840, 876],
        vec![77425, 2123, 2301, 1840, 1994, 2513, 1371, 874, 768, 549, 208, INF, 3, 707, 743],
        vec![77429, 2127, 2305, 1843, 1998, 2516, 1374, 877, 771, 553, 211, 3, INF, 711, 747],
        vec![78007, 2704, 2882, 2421, 2575, 3094, 1952, 1455, 1349, 1130, 789, 656, 660, INF, 117],
        vec![78094, 2792, 2969, 2508, 2662, 3181, 2039, 1542, 1436, 1218, 876, 743, 747, 117, INF],
    ];

    let max_iterations = 10_000_000;
    let start = Instant::now();
    let (best_cost, best_path) = branch_and_bound(graph.clone(), max_iterations);
    let elapsed = start.elapsed();

    println!("Time: {:?}", elapsed);
    println!("Minimum cost: {}", best_cost);
    println!("Best path: {:?}", best_path);
    println!("Calc path: {}", calculate_distance(&graph, &best_path));
}
