package rtree

import (
	"testing"
	"github.com/intdxdt/mbr"
	"github.com/franela/goblin"
	"time"
)

type Pnt struct {
	x, y float64
}

func (pt *Pnt) BBox() *mbr.MBR {
	return &mbr.MBR{pt.x, pt.y, pt.x + 2, pt.y + 2}
}

type nodeParent struct {
	wkt      string
	children []string
}

func printRtree(a *node) []*nodeParent {
	var tokens []*nodeParent
	if a == nil {
		return tokens
	}
	var nd *node
	var stack []*node
	stack = append(stack, a)
	for len(stack) > 0 {
		nd, stack = popNode(stack)
		var parent = &nodeParent{wkt: nd.bbox.String()}
		//adopt children on stack and let node go out of scope
		for i := range nd.children {
			if len(nd.children[i].children) > 0 {
				stack = append(stack, &nd.children[i])
				parent.children = append(parent.children, nd.children[i].bbox.String())
			}
		}

		if len(parent.children) > 0 {
			tokens = append(tokens, parent)
		}

	}
	return tokens
}

func TestRtree(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("rtree : node, leaf, inode", func() {
		var pt = &Pnt{0, 0}
		var item = pt
		var pth NodePath
		var b = newNode(item, 0, true, nil)

		pth = append(pth, b)
		pth = append(pth, b)
		pth = append(pth, b)

		var n = newNode(item, 1, false, pth)
		var items = make([]BoxObj, 0, 10)
		var nodes = make(NodePath, 0, 0)

		items = append(items, item)
		nodes = append(nodes, b)

		g.It("type node check ", func() {
			g.Assert(b.GetItem().(*Pnt)).Eql(pt)
			g.Assert(b.leaf).Equal(true)
			g.Assert(n.leaf).Equal(false)
			g.Assert(len(n.children)).Equal(3)

			g.Assert(b.height).Equal(0)
			g.Assert(n.height).Equal(1)
			g.Assert(len(b.children)).Equal(0)

			g.Assert(b.item).Equal(item)
			g.Assert(b.BBox()).Equal(&mbr.MBR{0, 0, 2, 2})
			g.Assert(n.BBox()).Equal(&mbr.MBR{0, 0, 2, 2})

		})

	})

	g.Describe("build rtree by bulkload and onebyone insert", func() {
		g.Timeout(1 * time.Hour)
		var data = []mbr.MBR{
			{30.74124324842736, 1.5394264094726768, 35.574749911400275, 8.754917282902216}, {7.381378714281472, 64.86180480586492, 19.198256264240655, 68.0987794848029}, {55.08436657887449, 73.66959671568338, 64.23351406225139, 77.82561878388375}, {60.0196123955198, 57.30699966964475, 74.83877476995968, 71.6532691707469}, {70.41627091953383, 51.438036044803454, 80.79446144066551, 55.724409001469795}, {6.483303127937942, 80.37332301675087, 6.50914529921677, 82.02059482017387}, {46.67649373819957, 64.24510021830747, 49.2622050275365, 72.2146377872009}, {13.896809634528902, 52.75698091860803, 27.3474212705194, 59.708006858014954}, {45.352809515933565, 67.57279878792961, 57.71107486286911, 80.63410132702094}, {58.12911437270156, 21.853066059797676, 72.6816258699198, 25.407156729750344}, {1.228055380876119, 71.76243208229317, 3.921389356330319, 71.81985676158466}, {24.176338710683243, 40.468612774474124, 30.028427694218617, 54.92587462821439}, {75.90272549273205, 70.78950716967577, 90.24958662679839, 73.14532201100896},
			{81.17621599681077, 43.92908059235767, 90.4623706429688, 45.683200269169774}, {10.765947677699227, 81.39085907882142, 16.395569791144023, 89.08943214908905}, {54.98460948258535, 75.98770610541906, 63.17175560560506, 89.58032814388704}, {42.47968070466303, 70.33863394618999, 53.969718678982176, 81.12499083427267}, {56.597735425362224, 22.872881616226724, 58.02513594712652, 29.461626653458254}, {28.072656807817236, 3.648771707777917, 32.25507880635046, 14.896956422497794}, {49.07401457054004, 65.43509168217955, 50.276686480083214, 72.13126764274583}, {66.92950379018822, 7.40714495221543, 78.79495207418685, 15.349551257658238}, {70.05814537971477, 81.30351958853318, 71.64399891813584, 91.16708488214654}, {21.4511094375575, 69.72891964401825, 31.722373869482286, 80.3256272486247}, {40.232777196706415, 26.760849136982976, 52.202812069867704, 34.21206366219117}, {2.368032571076858, 16.296113490306034, 12.33270360370716, 30.694571126478845}, {9.01855144170366, 55.970132314222134, 23.827554767436514, 60.48030769802354}, {80.61271599954031, 36.74002124278151, 91.79275857224492, 46.9506194268175}, {50.34990344969663, 81.49769656350821, 63.617315842569894, 83.30755417296216}, {39.18113381327339, 62.28148778267892, 46.4815234729194, 67.41798018502531},
			{29.998914416747247, 11.59250655284693, 33.376874697786775, 12.379204853229147}, {81.64879583058361, 25.545401825528394, 93.4343371235951, 37.16442658088167}, {38.58905494531754, 31.87745905840195, 41.7616624289497, 38.45823126735888}, {0.9178278426197698, 24.298283582889418, 13.300394793306303, 29.32894041204992}, {65.26849055356847, 81.26949067385523, 69.4019309878049, 95.14982799740329}, {41.57395146960945, 42.58630560128803, 44.74131455539111, 52.67240067840212}, {78.75491794797742, 24.519635432090283, 86.62303951191035, 27.152009252646756}, {57.413508019097335, 16.222132563535784, 64.52460425468645, 26.468580365950785}, {38.70624110521209, 63.6483778012707, 42.81587531412866, 76.69707330624905}, {45.79681150909137, 40.50191132346466, 56.183424730475984, 45.059343488954596}, {59.12908726623217, 61.8670788267583, 72.67061706103317, 74.71825120772677}, {53.530204647536515, 22.210826106446316, 56.19567351522378, 36.70783763707212}, {66.56685327399163, 41.84620000931149, 67.95502218856858, 51.90145172377749}, {13.647425280602949, 48.287305203367325, 14.605520880072303, 50.785335362500966}, {9.580714642281816, 71.82612512759374, 22.052586035203777, 78.60447881685704}, {42.52476287398914, 31.798014129805892, 47.30017532169579, 43.32042676277269}, {15.231406548475704, 20.91813524362627, 27.999049905750184, 33.12719299053375}, {68.25622304622375, 36.45344418954924, 75.12753345668115, 42.96412962336906},
			{24.674565636296396, 61.64103736035227, 33.35950737775334, 68.17273669513995}, {27.860994552259186, 54.07784655778231, 37.454370732019164, 55.03748662118532}, {12.989350409059881, 12.850601894307912, 19.63701743062105, 24.447201165885136}, {54.351699198645946, 38.669663277102835, 62.70698234918281, 50.77799147478973}, {5.195541592506005, 27.378150434771385, 12.470640457055284, 31.42600927621769}, {50.42859019394414, 76.74400020764121, 61.43712226636309, 81.94545584300995}, {78.94947703076804, 80.53231036050055, 80.65894334965007, 80.53525709875574}, {25.444253736005553, 7.68730085456098, 31.065085510940172, 20.3498357552189}, {67.23805308545823, 13.569845282055715, 72.08492158784647, 28.386336312117162}, {73.53304448250748, 72.95399805919209, 78.88399497592506, 86.10583870486123}, {5.128991214483967, 46.433989769953975, 10.301559209436643, 47.47697754635162}, {34.345971501358505, 37.67046253655506, 46.65109226249595, 43.20929547370596}, {46.288476425780644, 83.24699351224912, 53.04617705157806, 95.25275555638714}, {2.3371253741744717, 67.38540121025542, 13.258004924360035, 67.9350571187773}, {81.50701949936798, 12.96213890790966, 90.69810567341676, 26.897004984394016}, {19.618219504752606, 35.07620582977229, 22.719692101944606, 35.682818900087824}, {12.212115116661117, 56.27156067476181, 15.934817779897248, 62.75070913000411}, {68.37555295280667, 52.219237356472945, 68.38823378366567, 63.48647393313754},
			{30.62554452606222, 60.101485548798514, 37.063824618295754, 71.04525498101337}, {56.032005794131614, 71.80616208209968, 67.22546752158931, 83.70215276205255}, {20.14317265947747, 73.77798886182363, 34.25432987619779, 87.24104072094495}, {10.507678860183212, 66.06446404977234, 22.91945017863563, 73.50576752587352}, {26.0796380640738, 39.08543029877627, 37.497243272316375, 42.198598580655705}, {58.204665266130036, 58.20119021138755, 66.86094220293387, 61.613651791527374}, {40.43959914994069, 2.5737454435527933, 47.14440867190218, 10.136829689608973}, {81.61166337839565, 57.04686555019882, 82.13024015743876, 60.52557802686094}, {1.1438702774984308, 64.4390551345789, 1.207827079116793, 74.94606495692364}, {22.698477311365394, 31.694032934311718, 23.012351437738243, 34.826851291697004}, {58.23302290469934, 63.09245797513119, 63.89603555830784, 71.13299682623365}, {1.1209075169457285, 81.28342384198416, 2.010664217814431, 85.39246047317187}, {12.031894943077951, 47.03188640891187, 17.157531829906453, 58.84050109551066}, {25.175447117884868, 53.84501614745653, 29.018643250506607, 59.38873449198591}, {2.2848309030370015, 13.908167333298184, 9.169561431787841, 19.16049137202979}, {50.013550661499245, 78.5109200392331, 61.27884750099618, 90.82242857844415},
			{60.35181123067779, 50.30720879159393, 66.40423614499642, 62.711248070454005}, {12.818633233242565, 80.69085735063159, 25.51374909020891, 93.22537975149076}, {13.89435574446365, 30.374627423660982, 26.014177608552792, 40.22893652344269}, {68.59949104329682, 71.57717815724429, 71.14413101711249, 81.32143731631942}, {8.759053910523154, 40.17136447593845, 22.076247428918848, 51.97034411093291}, {75.0237223114521, 10.812195153356786, 75.45859644475163, 24.680056123348074}, {37.640987086884465, 44.31736944555115, 46.79079124130418, 52.298119297002756}, {77.86465045295246, 69.74685405122065, 91.0727578759392, 81.32602647164121}, {41.571023531510896, 41.188931957868, 47.81613155473583, 53.78551712929363}, {46.21623238891625, 12.566288400974617, 60.42998852835609, 23.520076065312416}, {39.651498265328506, 13.503482197678323, 50.2456922936693, 17.970333385957133}, {22.002987425318885, 4.223514231931571, 24.39665459195155, 17.79996696134728}, {10.238509846935935, 17.775671898372956, 24.90139389081459, 30.900047607940877}, {11.945673076143192, 11.005643838128806, 14.458677679728162, 25.935774067123525}, {34.15254570484473, 32.9087837466544, 39.806374568647804, 45.792474254223166}, {1.2619249479259986, 73.38259039620652, 5.732709854315865, 82.08100065666045},
			{68.88687814624431, 70.06499982957165, 70.86758866753506, 78.39070584782843}, {53.346140703038856, 38.61621943306142, 58.18001677406793, 46.227279405415416}, {60.91283806646173, 5.328797186659199, 70.97382774644399, 11.165367727083606},
		}

		var tree = NewRTree(9)
		var length = len(data)
		var data1By1 = data[:length:length]
		for i := range data1By1 {
			tree.Insert(&data1By1[i])
		}

		g.It("same root bounds for : bulkload & single insert ", func() {
			var nodeSize int
			var oneT = NewRTree(9)
			var oneDeft = NewRTree(nodeSize)
			var bulkT = NewRTree(9)

			//one by one
			var length = len(data)
			var dataOnebyone = data[:length:length]
			for i := range dataOnebyone {
				//fmt.Println(i, " -> ", len(oneT.Data.children))
				oneT.Insert(&dataOnebyone[i])
			}
			//fill zero size
			for i := range dataOnebyone {
				oneDeft.Insert(&dataOnebyone[i])
			}

			var oneMbr = oneT.Data.bbox
			var oneDefMbr = oneDeft.Data.bbox

			//fmt.Println(oneMbr.String())

			//bulkload
			var dataBulkload = data[:length:length]
			var bulkItems = make([]BoxObj, len(dataBulkload))
			for i := range bulkItems {
				bulkItems[i] = &dataBulkload[i]
			}
			bulkT.Load(bulkItems)
			bukMbr := bulkT.Data.bbox

			g.Assert(oneMbr).Eql(oneDefMbr)
			g.Assert(oneMbr).Eql(bukMbr)
			g.Assert(len(bulkT.Data.children)).Equal(len(oneT.Data.children))

			//var tokens = print_RTree(oneT.Data)
			//for _, tok := range tokens {
			//	fmt.Println(tok.wkt)
			//	for _, ch := range tok.children {
			//		fmt.Println("    " + ch)
			//	}
			//}
		})

	})

	g.Describe("build rtree by and remove all", func() {

		var data = []mbr.MBR{
			{30.74124324842736, 1.5394264094726768, 35.574749911400275, 8.754917282902216}, {7.381378714281472, 64.86180480586492, 19.198256264240655, 68.0987794848029}, {55.08436657887449, 73.66959671568338, 64.23351406225139, 77.82561878388375}, {60.0196123955198, 57.30699966964475, 74.83877476995968, 71.6532691707469}, {70.41627091953383, 51.438036044803454, 80.79446144066551, 55.724409001469795}, {6.483303127937942, 80.37332301675087, 6.50914529921677, 82.02059482017387}, {46.67649373819957, 64.24510021830747, 49.2622050275365, 72.2146377872009}, {13.896809634528902, 52.75698091860803, 27.3474212705194, 59.708006858014954}, {45.352809515933565, 67.57279878792961, 57.71107486286911, 80.63410132702094}, {58.12911437270156, 21.853066059797676, 72.6816258699198, 25.407156729750344}, {1.228055380876119, 71.76243208229317, 3.921389356330319, 71.81985676158466}, {24.176338710683243, 40.468612774474124, 30.028427694218617, 54.92587462821439}, {75.90272549273205, 70.78950716967577, 90.24958662679839, 73.14532201100896},
			{81.17621599681077, 43.92908059235767, 90.4623706429688, 45.683200269169774}, {10.765947677699227, 81.39085907882142, 16.395569791144023, 89.08943214908905}, {54.98460948258535, 75.98770610541906, 63.17175560560506, 89.58032814388704}, {42.47968070466303, 70.33863394618999, 53.969718678982176, 81.12499083427267}, {56.597735425362224, 22.872881616226724, 58.02513594712652, 29.461626653458254}, {28.072656807817236, 3.648771707777917, 32.25507880635046, 14.896956422497794}, {49.07401457054004, 65.43509168217955, 50.276686480083214, 72.13126764274583}, {66.92950379018822, 7.40714495221543, 78.79495207418685, 15.349551257658238}, {70.05814537971477, 81.30351958853318, 71.64399891813584, 91.16708488214654}, {21.4511094375575, 69.72891964401825, 31.722373869482286, 80.3256272486247}, {40.232777196706415, 26.760849136982976, 52.202812069867704, 34.21206366219117}, {2.368032571076858, 16.296113490306034, 12.33270360370716, 30.694571126478845}, {9.01855144170366, 55.970132314222134, 23.827554767436514, 60.48030769802354}, {80.61271599954031, 36.74002124278151, 91.79275857224492, 46.9506194268175}, {50.34990344969663, 81.49769656350821, 63.617315842569894, 83.30755417296216}, {39.18113381327339, 62.28148778267892, 46.4815234729194, 67.41798018502531},
			{29.998914416747247, 11.59250655284693, 33.376874697786775, 12.379204853229147}, {81.64879583058361, 25.545401825528394, 93.4343371235951, 37.16442658088167}, {38.58905494531754, 31.87745905840195, 41.7616624289497, 38.45823126735888}, {0.9178278426197698, 24.298283582889418, 13.300394793306303, 29.32894041204992}, {65.26849055356847, 81.26949067385523, 69.4019309878049, 95.14982799740329}, {41.57395146960945, 42.58630560128803, 44.74131455539111, 52.67240067840212}, {78.75491794797742, 24.519635432090283, 86.62303951191035, 27.152009252646756}, {57.413508019097335, 16.222132563535784, 64.52460425468645, 26.468580365950785}, {38.70624110521209, 63.6483778012707, 42.81587531412866, 76.69707330624905}, {45.79681150909137, 40.50191132346466, 56.183424730475984, 45.059343488954596}, {59.12908726623217, 61.8670788267583, 72.67061706103317, 74.71825120772677}, {53.530204647536515, 22.210826106446316, 56.19567351522378, 36.70783763707212}, {66.56685327399163, 41.84620000931149, 67.95502218856858, 51.90145172377749}, {13.647425280602949, 48.287305203367325, 14.605520880072303, 50.785335362500966}, {9.580714642281816, 71.82612512759374, 22.052586035203777, 78.60447881685704}, {42.52476287398914, 31.798014129805892, 47.30017532169579, 43.32042676277269}, {15.231406548475704, 20.91813524362627, 27.999049905750184, 33.12719299053375}, {68.25622304622375, 36.45344418954924, 75.12753345668115, 42.96412962336906},
			{24.674565636296396, 61.64103736035227, 33.35950737775334, 68.17273669513995}, {27.860994552259186, 54.07784655778231, 37.454370732019164, 55.03748662118532}, {12.989350409059881, 12.850601894307912, 19.63701743062105, 24.447201165885136}, {54.351699198645946, 38.669663277102835, 62.70698234918281, 50.77799147478973}, {5.195541592506005, 27.378150434771385, 12.470640457055284, 31.42600927621769}, {50.42859019394414, 76.74400020764121, 61.43712226636309, 81.94545584300995}, {78.94947703076804, 80.53231036050055, 80.65894334965007, 80.53525709875574}, {25.444253736005553, 7.68730085456098, 31.065085510940172, 20.3498357552189}, {67.23805308545823, 13.569845282055715, 72.08492158784647, 28.386336312117162}, {73.53304448250748, 72.95399805919209, 78.88399497592506, 86.10583870486123}, {5.128991214483967, 46.433989769953975, 10.301559209436643, 47.47697754635162}, {34.345971501358505, 37.67046253655506, 46.65109226249595, 43.20929547370596}, {46.288476425780644, 83.24699351224912, 53.04617705157806, 95.25275555638714}, {2.3371253741744717, 67.38540121025542, 13.258004924360035, 67.9350571187773}, {81.50701949936798, 12.96213890790966, 90.69810567341676, 26.897004984394016}, {19.618219504752606, 35.07620582977229, 22.719692101944606, 35.682818900087824}, {12.212115116661117, 56.27156067476181, 15.934817779897248, 62.75070913000411}, {68.37555295280667, 52.219237356472945, 68.38823378366567, 63.48647393313754},
			{30.62554452606222, 60.101485548798514, 37.063824618295754, 71.04525498101337}, {56.032005794131614, 71.80616208209968, 67.22546752158931, 83.70215276205255}, {20.14317265947747, 73.77798886182363, 34.25432987619779, 87.24104072094495}, {10.507678860183212, 66.06446404977234, 22.91945017863563, 73.50576752587352}, {26.0796380640738, 39.08543029877627, 37.497243272316375, 42.198598580655705}, {58.204665266130036, 58.20119021138755, 66.86094220293387, 61.613651791527374}, {40.43959914994069, 2.5737454435527933, 47.14440867190218, 10.136829689608973}, {81.61166337839565, 57.04686555019882, 82.13024015743876, 60.52557802686094}, {1.1438702774984308, 64.4390551345789, 1.207827079116793, 74.94606495692364}, {22.698477311365394, 31.694032934311718, 23.012351437738243, 34.826851291697004}, {58.23302290469934, 63.09245797513119, 63.89603555830784, 71.13299682623365}, {1.1209075169457285, 81.28342384198416, 2.010664217814431, 85.39246047317187}, {12.031894943077951, 47.03188640891187, 17.157531829906453, 58.84050109551066}, {25.175447117884868, 53.84501614745653, 29.018643250506607, 59.38873449198591}, {2.2848309030370015, 13.908167333298184, 9.169561431787841, 19.16049137202979}, {50.013550661499245, 78.5109200392331, 61.27884750099618, 90.82242857844415},
			{60.35181123067779, 50.30720879159393, 66.40423614499642, 62.711248070454005}, {12.818633233242565, 80.69085735063159, 25.51374909020891, 93.22537975149076}, {13.89435574446365, 30.374627423660982, 26.014177608552792, 40.22893652344269}, {68.59949104329682, 71.57717815724429, 71.14413101711249, 81.32143731631942}, {8.759053910523154, 40.17136447593845, 22.076247428918848, 51.97034411093291}, {75.0237223114521, 10.812195153356786, 75.45859644475163, 24.680056123348074}, {37.640987086884465, 44.31736944555115, 46.79079124130418, 52.298119297002756}, {77.86465045295246, 69.74685405122065, 91.0727578759392, 81.32602647164121}, {41.571023531510896, 41.188931957868, 47.81613155473583, 53.78551712929363}, {46.21623238891625, 12.566288400974617, 60.42998852835609, 23.520076065312416}, {39.651498265328506, 13.503482197678323, 50.2456922936693, 17.970333385957133}, {22.002987425318885, 4.223514231931571, 24.39665459195155, 17.79996696134728}, {10.238509846935935, 17.775671898372956, 24.90139389081459, 30.900047607940877}, {11.945673076143192, 11.005643838128806, 14.458677679728162, 25.935774067123525}, {34.15254570484473, 32.9087837466544, 39.806374568647804, 45.792474254223166}, {1.2619249479259986, 73.38259039620652, 5.732709854315865, 82.08100065666045},
			{68.88687814624431, 70.06499982957165, 70.86758866753506, 78.39070584782843}, {53.346140703038856, 38.61621943306142, 58.18001677406793, 46.227279405415416}, {60.91283806646173, 5.328797186659199, 70.97382774644399, 11.165367727083606},
		}
		var query = mbr.MBR{0, 0, 100, 100}
		var tree = NewRTree()
		var length = len(data)
		var dataOnebyone = data[:length:length]
		for i := range dataOnebyone {
			tree.Insert(&dataOnebyone[i])
		}

		g.It("same root bounds for : bulkload & single insert ", func() {
			var res = tree.Search(query)
			for i := range res {
				tree.RemoveObj(res[i])
			}
			g.Assert(tree.IsEmpty()).IsTrue()
			g.Assert(len(tree.Data.children)).Equal(0)
			g.Assert(tree.Data.bbox).Eql(emptyMBR())

		})
	})

	g.Describe("search for items in tree", func() {

		var data = []mbr.MBR{
			{30.74124324842736, 1.5394264094726768, 35.574749911400275, 8.754917282902216}, {7.381378714281472, 64.86180480586492, 19.198256264240655, 68.0987794848029}, {55.08436657887449, 73.66959671568338, 64.23351406225139, 77.82561878388375}, {60.0196123955198, 57.30699966964475, 74.83877476995968, 71.6532691707469}, {70.41627091953383, 51.438036044803454, 80.79446144066551, 55.724409001469795}, {6.483303127937942, 80.37332301675087, 6.50914529921677, 82.02059482017387}, {46.67649373819957, 64.24510021830747, 49.2622050275365, 72.2146377872009}, {13.896809634528902, 52.75698091860803, 27.3474212705194, 59.708006858014954}, {45.352809515933565, 67.57279878792961, 57.71107486286911, 80.63410132702094}, {58.12911437270156, 21.853066059797676, 72.6816258699198, 25.407156729750344}, {1.228055380876119, 71.76243208229317, 3.921389356330319, 71.81985676158466}, {24.176338710683243, 40.468612774474124, 30.028427694218617, 54.92587462821439}, {75.90272549273205, 70.78950716967577, 90.24958662679839, 73.14532201100896},
			{81.17621599681077, 43.92908059235767, 90.4623706429688, 45.683200269169774}, {10.765947677699227, 81.39085907882142, 16.395569791144023, 89.08943214908905}, {54.98460948258535, 75.98770610541906, 63.17175560560506, 89.58032814388704}, {42.47968070466303, 70.33863394618999, 53.969718678982176, 81.12499083427267}, {56.597735425362224, 22.872881616226724, 58.02513594712652, 29.461626653458254}, {28.072656807817236, 3.648771707777917, 32.25507880635046, 14.896956422497794}, {49.07401457054004, 65.43509168217955, 50.276686480083214, 72.13126764274583}, {66.92950379018822, 7.40714495221543, 78.79495207418685, 15.349551257658238}, {70.05814537971477, 81.30351958853318, 71.64399891813584, 91.16708488214654}, {21.4511094375575, 69.72891964401825, 31.722373869482286, 80.3256272486247}, {40.232777196706415, 26.760849136982976, 52.202812069867704, 34.21206366219117}, {2.368032571076858, 16.296113490306034, 12.33270360370716, 30.694571126478845}, {9.01855144170366, 55.970132314222134, 23.827554767436514, 60.48030769802354}, {80.61271599954031, 36.74002124278151, 91.79275857224492, 46.9506194268175}, {50.34990344969663, 81.49769656350821, 63.617315842569894, 83.30755417296216}, {39.18113381327339, 62.28148778267892, 46.4815234729194, 67.41798018502531},
			{29.998914416747247, 11.59250655284693, 33.376874697786775, 12.379204853229147}, {81.64879583058361, 25.545401825528394, 93.4343371235951, 37.16442658088167}, {38.58905494531754, 31.87745905840195, 41.7616624289497, 38.45823126735888}, {0.9178278426197698, 24.298283582889418, 13.300394793306303, 29.32894041204992}, {65.26849055356847, 81.26949067385523, 69.4019309878049, 95.14982799740329}, {41.57395146960945, 42.58630560128803, 44.74131455539111, 52.67240067840212}, {78.75491794797742, 24.519635432090283, 86.62303951191035, 27.152009252646756}, {57.413508019097335, 16.222132563535784, 64.52460425468645, 26.468580365950785}, {38.70624110521209, 63.6483778012707, 42.81587531412866, 76.69707330624905}, {45.79681150909137, 40.50191132346466, 56.183424730475984, 45.059343488954596}, {59.12908726623217, 61.8670788267583, 72.67061706103317, 74.71825120772677}, {53.530204647536515, 22.210826106446316, 56.19567351522378, 36.70783763707212}, {66.56685327399163, 41.84620000931149, 67.95502218856858, 51.90145172377749}, {13.647425280602949, 48.287305203367325, 14.605520880072303, 50.785335362500966}, {9.580714642281816, 71.82612512759374, 22.052586035203777, 78.60447881685704}, {42.52476287398914, 31.798014129805892, 47.30017532169579, 43.32042676277269}, {15.231406548475704, 20.91813524362627, 27.999049905750184, 33.12719299053375}, {68.25622304622375, 36.45344418954924, 75.12753345668115, 42.96412962336906},
			{24.674565636296396, 61.64103736035227, 33.35950737775334, 68.17273669513995}, {27.860994552259186, 54.07784655778231, 37.454370732019164, 55.03748662118532}, {12.989350409059881, 12.850601894307912, 19.63701743062105, 24.447201165885136}, {54.351699198645946, 38.669663277102835, 62.70698234918281, 50.77799147478973}, {5.195541592506005, 27.378150434771385, 12.470640457055284, 31.42600927621769}, {50.42859019394414, 76.74400020764121, 61.43712226636309, 81.94545584300995}, {78.94947703076804, 80.53231036050055, 80.65894334965007, 80.53525709875574}, {25.444253736005553, 7.68730085456098, 31.065085510940172, 20.3498357552189}, {67.23805308545823, 13.569845282055715, 72.08492158784647, 28.386336312117162}, {73.53304448250748, 72.95399805919209, 78.88399497592506, 86.10583870486123}, {5.128991214483967, 46.433989769953975, 10.301559209436643, 47.47697754635162}, {34.345971501358505, 37.67046253655506, 46.65109226249595, 43.20929547370596}, {46.288476425780644, 83.24699351224912, 53.04617705157806, 95.25275555638714}, {2.3371253741744717, 67.38540121025542, 13.258004924360035, 67.9350571187773}, {81.50701949936798, 12.96213890790966, 90.69810567341676, 26.897004984394016}, {19.618219504752606, 35.07620582977229, 22.719692101944606, 35.682818900087824}, {12.212115116661117, 56.27156067476181, 15.934817779897248, 62.75070913000411}, {68.37555295280667, 52.219237356472945, 68.38823378366567, 63.48647393313754},
			{30.62554452606222, 60.101485548798514, 37.063824618295754, 71.04525498101337}, {56.032005794131614, 71.80616208209968, 67.22546752158931, 83.70215276205255}, {20.14317265947747, 73.77798886182363, 34.25432987619779, 87.24104072094495}, {10.507678860183212, 66.06446404977234, 22.91945017863563, 73.50576752587352}, {26.0796380640738, 39.08543029877627, 37.497243272316375, 42.198598580655705}, {58.204665266130036, 58.20119021138755, 66.86094220293387, 61.613651791527374}, {40.43959914994069, 2.5737454435527933, 47.14440867190218, 10.136829689608973}, {81.61166337839565, 57.04686555019882, 82.13024015743876, 60.52557802686094}, {1.1438702774984308, 64.4390551345789, 1.207827079116793, 74.94606495692364}, {22.698477311365394, 31.694032934311718, 23.012351437738243, 34.826851291697004}, {58.23302290469934, 63.09245797513119, 63.89603555830784, 71.13299682623365}, {1.1209075169457285, 81.28342384198416, 2.010664217814431, 85.39246047317187}, {12.031894943077951, 47.03188640891187, 17.157531829906453, 58.84050109551066}, {25.175447117884868, 53.84501614745653, 29.018643250506607, 59.38873449198591}, {2.2848309030370015, 13.908167333298184, 9.169561431787841, 19.16049137202979}, {50.013550661499245, 78.5109200392331, 61.27884750099618, 90.82242857844415},
			{60.35181123067779, 50.30720879159393, 66.40423614499642, 62.711248070454005}, {12.818633233242565, 80.69085735063159, 25.51374909020891, 93.22537975149076}, {13.89435574446365, 30.374627423660982, 26.014177608552792, 40.22893652344269}, {68.59949104329682, 71.57717815724429, 71.14413101711249, 81.32143731631942}, {8.759053910523154, 40.17136447593845, 22.076247428918848, 51.97034411093291}, {75.0237223114521, 10.812195153356786, 75.45859644475163, 24.680056123348074}, {37.640987086884465, 44.31736944555115, 46.79079124130418, 52.298119297002756}, {77.86465045295246, 69.74685405122065, 91.0727578759392, 81.32602647164121}, {41.571023531510896, 41.188931957868, 47.81613155473583, 53.78551712929363}, {46.21623238891625, 12.566288400974617, 60.42998852835609, 23.520076065312416}, {39.651498265328506, 13.503482197678323, 50.2456922936693, 17.970333385957133}, {22.002987425318885, 4.223514231931571, 24.39665459195155, 17.79996696134728}, {10.238509846935935, 17.775671898372956, 24.90139389081459, 30.900047607940877}, {11.945673076143192, 11.005643838128806, 14.458677679728162, 25.935774067123525}, {34.15254570484473, 32.9087837466544, 39.806374568647804, 45.792474254223166}, {1.2619249479259986, 73.38259039620652, 5.732709854315865, 82.08100065666045},
			{68.88687814624431, 70.06499982957165, 70.86758866753506, 78.39070584782843}, {53.346140703038856, 38.61621943306142, 58.18001677406793, 46.227279405415416}, {60.91283806646173, 5.328797186659199, 70.97382774644399, 11.165367727083606},
		}

		//queries
		//nothing
		var query1 = mbr.MBR{81.59858271428983, 88.95212575682031, 87.00714129337072, 92.42905627194374}
		var query2 = mbr.MBR{82.17807113347706, 83.15724156494792, 87.39346690616222, 84.70254401611389}
		var query3 = mbr.MBR{84.10969919743454, 72.14696160039038, 86.23449006778775, 79.10082263063724}
		var query4 = mbr.MBR{21.298871774427138, 1.1709155631470283, 36.23985259304277, 20.747325333798532}
		var query5 = mbr.MBR{0., 0., 100, 100}
		var query6 = mbr.MBR{182.17619056720642, 15.748541593521262, 205.43811579298725, 65.97783146157896}

		var tree = NewRTree()
		var bulkTree = NewRTree()

		var length = len(data)
		var dataOnebyone = data[:length:length]
		var dataBulkload = data[:length:length]
		var bulkItems = make([]BoxObj, len(dataBulkload))

		for i := range bulkItems {
			bulkItems[i] = &dataBulkload[i]
		}

		for i := range dataOnebyone {
			tree.Insert(&dataOnebyone[i])
		}
		bulkTree.Load(bulkItems)

		g.It("should return items found - one-by-one", func() {
			var res1 = tree.Search(query1)
			var res2 = tree.Search(query2)
			var res3 = tree.Search(query3)
			var res4 = tree.Search(query4)
			var res5 = tree.Search(query5)
			var res6 = tree.Search(query6)

			g.Assert(len(res1)).Equal(0)
			g.Assert(len(res2)).Equal(0)
			g.Assert(len(res3)).Equal(2)
			g.Assert(len(res4)).Equal(6)
			g.Assert(len(res5)).Equal(len(data))
			g.Assert(len(res6)).Equal(0)
			g.Assert(len(tree.All())).Equal(len(data))
		})

		g.It("should return items found - bulk loaded tree", func() {
			var res1 = bulkTree.Search(query1)
			var res2 = bulkTree.Search(query2)
			var res3 = bulkTree.Search(query3)
			var res4 = bulkTree.Search(query4)
			var res5 = bulkTree.Search(query5)
			var res6 = bulkTree.Search(query6)

			g.Assert(len(res1)).Equal(0)
			g.Assert(len(res2)).Equal(0)
			g.Assert(len(res3)).Equal(2)
			g.Assert(len(res4)).Equal(6)
			g.Assert(len(res5)).Equal(len(data))
			g.Assert(len(res6)).Equal(0)
			g.Assert(len(tree.All())).Equal(len(data))
		})
	})
}
