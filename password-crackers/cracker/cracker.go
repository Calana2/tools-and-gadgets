package cracker

import (
	"fmt"
	"sync"
)

type Combiner struct {
 dictionary string
 MinLength uint8
 MaxLength uint8
}

const ASCII_printable_characters = "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"

/* Creates a new Combiner.

 The var dictionary is the alphabet used by the combinator.
 If dictionary is "", then the 94 printable ASCII characters are used.

 M is the minimum word length generated.

 N is the maximum word length generated.

*/
func NewCombiner(dictionary string, M, N uint8) (*Combiner, error) {
  if M == 0 || N == 0 { 
   return nil, fmt.Errorf("Word length can't be zero.")
  }

  if dictionary == "" {
   dictionary = ASCII_printable_characters
  }

  c := &Combiner{dictionary: dictionary, MinLength: M, MaxLength: N}
  return c,nil
}

func (c *Combiner) SetMinLength(value uint8) error {
  if value == 0 {  
   return fmt.Errorf("Word length can't be zero.")
  }
  c.MinLength = value
  return nil
}

func (c *Combiner) SetMaxLength(value uint8) error {
  if value == 0 {  
   return fmt.Errorf("Word length can't be zero.")
  }
  c.MaxLength = value
  return nil
}

// Generate words of M <= length <= N to a string channel 
// Ends with sync.WaitGroup.Done() 
func (c *Combiner) GenerateToPipe(pipe chan string, s *sync.WaitGroup) {
  for i := c.MinLength ; i <= c.MaxLength; i++ {
   c.generateCombinations(&c.dictionary,i,"",pipe)
  }
  s.Done()
} 

// recursive function used by GenerateToPipe
func (c *Combiner) generateCombinations(dictionary *string, length uint8, currentCombination string, pipe chan string) {
  if length == 0 {
   pipe <- currentCombination
   return
  }
  for _,char := range(*dictionary) {
    c.generateCombinations(dictionary,length-1,currentCombination+string(char), pipe)
  }
}

