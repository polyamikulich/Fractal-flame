## Некоторые рекомендации по генерации изображений

В качестве основы дефолтных афинных преобразований используются улучшенные преобразования для получения треугольника Серпинского. В качестве дефолтной функции используется линейная вариация с весом `1`. Напомню, что стандартный треугольник Серпинского получается путём применения следующих афинных преобразований:

```bash
go run ./cmd/fractalflame --affine-params 0.5,0,0,0,0.5,0/0.5,0,0.5,0,0.5,0/0.5,0,0,0,0.5,0.5 -i 5000000 -o doc_ex/image/result1.png
```

(см. `doc_ex/image/result1.png`)

Эту и дальнейшие команды можно смело вставлять в консоль и запускать, дабы убедиться в том, что все работает. Однако можно и просто заглянуть с папку `doc_ex/image` и посмотреть сгенерированные картинки.

В общем виде для треугольника Серпинского подойдут любые преобразования вида:

```
k,0,0,0,k,0/k,0,0,0,k,l/k,0,l,0,k,0
```

```bash
go run ./cmd/fractalflame --affine-params 0.6,0,0,0,0.6,0/0.6,0,0.4,0,0.6,0/0.6,0,0,0,0.6,0.4 -i 5000000 -o doc_ex/image/result2.png
```

(см. `doc_ex/image/result2.png`)

Можно использовать также следующие базовые преобразования:

1. Поворот на угол 

`(Cos,Sin,0,-Sin,Cos,0)`

Так, например, желая повернуть треугольник Серпинского на +-30 градусов, можно использовать следующую команду:

```bash
go run ./cmd/fractalflame --affine-params 0.5,0,0,0,0.5,0/0.5,0,0.5,0,0.5,0/0.5,0,0,0,0.5,0.5/0.866,0.5,0,-0.5,0.866,0/0.866,-0.5,0,0.5,0.866,0 -i 5000000 -o doc_ex/image/result3.png
```

(см. `doc_ex/image/result3.png`)

2. Сдвиг в точку с координатами (х,у)

`(k,0,х,0,k,у)`

```bash
go run ./cmd/fractalflame --affine-params 0.5,0,0,0,0.5,0/0.5,0,0.5,0,0.5,0/0.5,0,0,0,0.5,0.5/0.5,0,-0.5,0,0.5,0 -i 5000000 -o doc_ex/image/result4.png
```

(см. `doc_ex/image/result4.png`)

3. Растяжение в k раз по оси Х

`(k,0,0,0,1,0)`

4. Растяжение в n раз по оси Y

`(1,0,0,0,k,0)`

```bash
go run ./cmd/fractalflame --affine-params 0.5,0,0,0,0.5,0/0.5,0,0.5,0,0.5,0/0.5,0,0,0,0.5,0.5/2,0,0,0,0.5,0/0.5,0,0,0,-5,0 -i 5000000 -o doc_ex/image/result5.png
```

(см. `doc_ex/image/result5.png`)

В качестве дефолтных преобразований используется следующий набор афинных преобразований:

```
0.5,0,0,0,0.5,0
0.5,0,0.5,0,0.5,0
0.5,0,-0.5,0,0.5,0
0.5,0,0,0,0.5,0.5
0.5,0,0,0,0.5,-0.5
```

```bash
go run ./cmd/fractalflame -i 5000000 -o doc_ex/image/result.png
```

(см. `doc_ex/image/result.png`)

Поверх можно применять уже всякие разные другие вариации:

```bash
go run ./cmd/fractalflame -i 5000000 -o doc_ex/image/result_swirl.png -f swirl:1
```

```bash
go run ./cmd/fractalflame -i 5000000 -o doc_ex/image/result_horseshoe.png -f horseshoe:1
```

```bash
go run ./cmd/fractalflame -i 5000000 -o doc_ex/image/result_polar.png -f polar:1
```

```bash
go run ./cmd/fractalflame -i 5000000 -o doc_ex/image/result_spiral.png -f spiral:1 --gamma-correction 1
```

```bash
go run ./cmd/fractalflame --affine-params 0.5,0,0,0,0.5,0/0.5,0,0.5,0,0.5,0/0.5,0,-0.5,0,0.5,0/0.5,0,0,0,0.5,0.5/0.5,0,0,0,0.5,-0.5/100,0,0,0,1,0/1,0,0,0,0.1,0 -f heart:0.4 -i 5000000 -o doc_ex/image/result_heart.png --seed 78 --gamma-correction 1
```

Помимо гамма-коррекции можно также воспользоваться симметрией:

```bash
go run ./cmd/fractalflame -i 5000000 -o doc_ex/image/result_symmetry3.png -s 3
```

