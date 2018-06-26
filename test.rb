def block_id(r, c)
  ((r / 3) * 3) + (c / 3)
end

(0...9).each do |r|
  (0...9).each do |c|
    printf("%s ", block_id(r, c))
  end
  puts
end

(0 ... 9).each do |b|
  (0 ... 9).each do |i|
    offset_r = b / 3 * 3
    offset_c = b % 3 * 3
    r = (i / 3) + offset_r
    c = (i % 3) + offset_c
    printf("(%s, %s, %s)", b, r, c)
  end
  puts
end

(0 .. 9).each do |i|
  printf("0x%03x  // %s %09b\n", 1 << i-1, i, 1 << i-1)
end
