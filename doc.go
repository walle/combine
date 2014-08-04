// Package combine contains functionality for combining templates to files by reading other files.
//
// The main parts are the interfaces Includer and Combiner.
//
// If you wish to extend the library with for an example another template package (TemplateIncluder uses text/template)
// implement the methods of Includer and use the defined combiners.
//
// If you wish to extent the library with another combiner, look at YuiMinifyCombiner and see how it is used
// to extend the standard combiner with aditional functionality YuiMinifyCombiner#Combine.
package combine
